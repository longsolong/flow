// Copyright 2017-2019, Square, Inc.

package traverser

import (
	"context"
	"fmt"
	"github.com/longsolong/flow/pkg/infra"
	"github.com/longsolong/flow/pkg/orchestration/job"
	"github.com/longsolong/flow/pkg/orchestration/job_runner/traverser"
	"github.com/longsolong/flow/pkg/orchestration/job_runner/runner"
	"github.com/longsolong/flow/pkg/orchestration/single_processor/graph"
	"github.com/longsolong/flow/pkg/workflow/state"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
	"sync/atomic"
	"time"
)

// Traverser ...
type Traverser struct {
	grapher *graph.Grapher

	reaper traverser.JobReaper

	runJobChan  chan job.Job  // jobs to be run
	doneJobChan chan job.Job  // jobs that are done
	doneChan    chan struct{} // closed when traverser finishes running

	stopMux     *sync.RWMutex // lock around checks to stopped
	stopped     bool          // has traverser been stopped
	stopChan    chan struct{} // don't run jobs in runJobs
	pendingChan chan struct{} // runJobs closes on return
	pending     int64         // N runJob goroutines are pending runnerRepo.Set

	runnerRepo runner.Repo // stores actively running jobs

	logger *infra.Logger

	stopTimeout time.Duration // Time to wait for jobs to stop
	sendTimeout time.Duration // Time to wait for a job to send on doneJobChan.
}

// NewTraverser ...
func NewTraverser(grapher *graph.Grapher, logger *infra.Logger, stopTimeout, sendTimeout time.Duration) *Traverser {
	// Channels used to communicate between traverser + reaper(s)
	doneJobChan := make(chan job.Job)
	runJobChan := make(chan job.Job)

	// Each traverser has its own runner repo because it's keyed on job ID and
	// job IDs are unique per-chain, not globally.
	runnerRepo := runner.NewRepo()

	return &Traverser{
		grapher: grapher,

		runnerRepo: runnerRepo,

		runJobChan:  runJobChan,
		doneJobChan: doneJobChan,

		doneChan:    make(chan struct{}),
		stopChan:    make(chan struct{}),
		pendingChan: make(chan struct{}),

		logger: logger,

		stopMux:     &sync.RWMutex{},
		stopTimeout: stopTimeout,
		sendTimeout: sendTimeout,
	}
}

// Run runs all jobs in the chain and blocks until the chain finishes running, is
// stopped.
func (t *Traverser) Run(ctx context.Context) {
	logger := t.logger.Log
	logger.Info("traverser.Run call")
	defer logger.Info("traverser.Run return")

	// Start a goroutine to run jobs. This consumes runJobChan. When jobs are done,
	// they're sent to doneJobChan, which a reaper consumes. This goroutine returns
	// when runJobChan is closed below.
	go t.runJobs(ctx)

	// Enqueue all the first runnable jobs
	for _, j := range t.grapher.Chain.RunnableJobs() {
		node := t.grapher.Chain.Vertices[j.StepID()]
		fields := []zapcore.Field{
			zap.String("job_id", j.StepID().String()),
			zap.String("job_name", node.Name),
		}
		logger.Info("initial job", fields...)
		t.runJobChan <- *j
	}

	// Start a goroutine to reap done jobs. The runningReaper consumes from
	// doneJobChan and sends the next jobs to be run to runJobChan. Stop()
	// calls t.reaper.Stop(), which is this reaper. The close(t.runJobChan)
	// causes runJobs() (started above ^) to return.
	runningReaperChan := make(chan struct{})
	t.reaper = NewRunningChainReaper(t.grapher, t.logger, t.doneJobChan, t.runJobChan)
	go func() {
		defer close(runningReaperChan) // indicate reaper is done (see select below)
		defer close(t.runJobChan)      // stop runJobs goroutine
		t.reaper.Run(ctx)
	}()

	// Wait for running reaper to be done.
	select {
	case <-runningReaperChan:
		// If running reaper is done because traverser was stopped, we will
		// wait for Stop() to finish. Otherwise, the chain finished normally
		// (completed or failed) and we can return right away.
		t.stopMux.Lock()
		if !t.stopped {
			t.stopMux.Unlock()
			return
		}
		t.stopMux.Unlock()
	}

	// Traverser is being stopped - wait for that to finish before
	// returning.
	select {
	case <-t.doneChan:
		// Stopped successfully - nothing left to do.
		return
	case <-time.After(20 * time.Second):
		// Failed to stop in a reasonable amount of time.
		// Log the failure and return.
		logger.Warn("stopping the job chain took too long. Exiting...")
		return
	}
}

// Stop stops the running job chain and stopping all currently running jobs.
func (t *Traverser) Stop(ctx context.Context) error {
	// Don't do anything if the traverser has already been stopped.
	t.stopMux.Lock()
	defer t.stopMux.Unlock()
	if t.stopped {
		return nil
	}
	logger := t.logger.Log

	close(t.stopChan)
	t.stopped = true
	logger.Info("stopping traverser and all jobs")

	// Stop the runningReaper
	t.reaper.Stop(ctx) // blocks until runningReaper stops

	// Stop all job runners in the runner repo.
	timeout := time.After(t.stopTimeout)
	err := t.stopRunningJobs(ctx, timeout)
	if err != nil {
		// Don't return the error yet - we still want to wait for the stop
		// reaper to be done.
		err = fmt.Errorf("traverser was stopped, but encountered an error in the process: %s", err)
	}

	close(t.doneChan)
	return err
}

// runJobs loops on the runJobChan, and runs each job that comes through the
// channel. When the job is done, it sends the job out through the doneJobChan
// which is being consumed by a reaper.
func (t *Traverser) runJobs(ctx context.Context) {
	logger := t.logger.Log
	logger.Info("runJobs call")
	defer logger.Info("runJobs return")
	defer close(t.pendingChan)

	// Run all jobs that come in on runJobChan. The loop exits when runJobChan
	// is closed in the runningReaper goroutine in Run().
	for j := range t.runJobChan {
		// Don't run the job if traverser stopped. In this case,
		// drain runJobChan to prevent runningReaper from blocking (the chan is
		// unbuffered).
		//
		// Must check before running goroutine because Run() closes runJobChan
		// when the runningReaper is done. Then this loop will end and close
		// pendingChan which stopRunningJobs blocks on. Since this check happens
		// in loop not goroutine, a closed pendingChan means it's been checked
		// for all jobs and either the job did not run or it did with pending+1
		// because the loop won't finish until running all code before the goroutine
		// is launched.
		select {
		case <-t.stopChan:
			fields := []zapcore.Field{
				zap.String("job_id", j.StepID().String()),
			}
			logger.Info("not running job %s: traverser stopped", fields...)
			continue
		default:
		}

		// Signal to stopRunningJobs that there's +1 goroutine
		atomic.AddInt64(&t.pending, 1)

		// Explicitly pass the job into the func, or all goroutines would share
		// the same loop "j" variable.
		go func(j job.Job) {
			fields := []zapcore.Field{
				zap.String("job_id", j.StepID().String()),
				zap.String("sequence_id", t.grapher.DAG.Vertices[j.StepID()].SequenceID.String()),
				zap.Int("sequence_try", int(t.grapher.DAG.Vertices[j.StepID()].SequenceRetry)),
			}

			// Always send the finished job to doneJobChan to be reaped. If the
			// reaper isn't reaping any more jobs (if this job took too long to
			// finish after being stopped), sending to doneJobChan won't be
			// possible - timeout after a while so we don't leak this goroutine.
			defer func() {
				select {
				case t.doneJobChan <- j: // reap the done job
				case <-time.After(t.sendTimeout):
					logger.Warn("timed out sending job to doneJobChan", fields...)
				}
				// Remove the job's runner from the repo (if it was ever added)
				// AFTER sending it to doneJobChan. This avoids a race condition
				// when the stopped + suspended reapers check if the runnerRepo
				// is empty.
				t.runnerRepo.Remove(j.StepID().String())
			}()

			// Increment sequence try count if this is sequence start job, which
			// currently means sequenceId == job.Id.
			if t.grapher.Chain.IsSequenceStartJob(j.StepID()) {
				t.grapher.Chain.IncrementSequenceTries(j.StepID(), 1)
				tryFields := append([]zapcore.Field(nil), fields...)
				tryFields = append(tryFields, zap.Uint("current", t.grapher.Chain.SequenceTries(j.StepID())))
				logger.Info("sequence try", tryFields...)
			}

			totalTries := t.grapher.Chain.JobTries(j.StepID())

			node := t.grapher.Chain.Vertices[j.StepID()]
			jobRunner := runner.NewRunner(j, t.grapher.Req, totalTries, node.Name, node.Retry, node.RetryWait, t.logger)

			// Add the runner to the repo. Runners in the repo are used
			// by the Stop methods on the traverser.
			// Then decrement pending to signal to stopRunningJobs that
			// there's one less goroutine it needs to wait for.
			t.runnerRepo.Set(j.StepID().String(), jobRunner)
			atomic.AddInt64(&t.pending, -1)

			// Run the job. This is a blocking operation that could take a long time.
			logger.Info("running job", fields...)
			t.grapher.Chain.SetJobState(j.StepID(), state.StateRunning)
			ret := jobRunner.Run(ctx)
			runFields := append([]zapcore.Field(nil), fields...)
			runFields = append(runFields, zap.String("state", state.StateText[ret.AtomReturn.State]))
			logger.Info("job done", runFields...)

			// We don't pass the Chain to the job runner, so it can't call this
			// itself. Instead, it returns how many tries it did, and we set it.
			t.grapher.Chain.IncrementJobTries(j.StepID(), ret.Tries)

			// Set job final state because this job is about to be reaped on
			// the doneJobChan, sent in this goroutine's defer func at top ^.
			j.State = ret.AtomReturn.State
		}(j)
	}
}

// stopRunningJobs stops all currently running jobs.
func (t *Traverser) stopRunningJobs(ctx context.Context, timeout <-chan time.Time) error {
	// To stop all running jobs without race conditions, we need to know:
	//   1. runJobs is done, won't start any more goroutines
	//   2. All in-flight runJob goroutines have added themselves to runner repo
	// First is easy: wait for it to close pendingChan. Second is like a wait
	// group wait: runJobs add +1 to pending when goroutine starts, and -1 after
	// it adds itself to runner repo. So all runJob goroutines have added
	// themselves to the runner repo when pending == 0.
	//
	// The shutdown sequence is:
	//   1. close(stopChan): runJob goroutines (RGs) don't run if closed. It's
	//      as if the job never ran. This allows runJobChan to drain and prevents
	//      runningReaper from blocking because the chan is unbuffered. This is done
	//      in the for loop, before launching the goroutine, so that a closed
	//      pendingChan (step 4) guarantees that runJobs either didn't run an
	//      RG or it did and added it to pending count (because the loop won't
	//      exit until running pre-goroutine code, and pendingChan is only closed
	//      after loop exits).
	//   2. Stop runningReaper: This stops new/next jobs into runJobChan, which
	//      is being drained because of step 1.
	//   3. close(runJobChan): When runningReaper.Run returns, the goroutine in
	//      in Traverser.Run closes runJobChan. Since runningReaper is only thing
	//      that sends to runJobChan, it must be closed like this so runningReaper
	//      doesn't panic on "send on closed channel".
	//   4. close(pendingChan): Given step 3 and step 1, eventually runJobChan
	//      will drain and runJobs() will return, closing pendingChan when it does.
	//   5. Call stopRunningJobs: This func waits for step 4, which ensures no
	//      more RGs. And given step 1, we're assured that all in-flight RGs
	//      have added themsevs to pending count. Therefore, this func waits for
	//      pending count == 0 which means all RGs have added themselves to the
	//      runner repo.
	//   6. Stop all active runners in runner repo.

	// Wait for runJobs to return
	select {
	case <-t.pendingChan:
	case <-timeout:
		return fmt.Errorf("stopRunningJobs: timeout waiting for pendingChan")
	}

	// Wait for in-flight runJob goroutines to add themselves to runner repo
	if n := atomic.LoadInt64(&t.pending); n > 0 {
		for atomic.LoadInt64(&t.pending) > 0 {
			select {
			case <-timeout:
				return fmt.Errorf("stopRunningJobs: timeout waiting for pending count")
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}

	logger := t.logger.Log
	// Stop all runners in parallel in case some jobs don't stop quickly
	activeRunners := t.runnerRepo.Items()
	fields := []zapcore.Field{
		zap.Int("num", len(activeRunners)),
	}
	logger.Info("stopping active job runners", fields...)
	var wg sync.WaitGroup
	hadError := false
	for jobID, activeRunner := range activeRunners {
		jobID := jobID
		wg.Add(1)
		go func(runner runner.Runner) {
			defer wg.Done()
			if err := runner.Stop(ctx); err != nil {
				fields := []zapcore.Field{
					zap.String("job_id", jobID),
					zap.String("error", err.Error()),
				}
				logger.Error("problem stopping job runner", fields...)
				hadError = true
			}
		}(activeRunner)
	}
	wg.Wait()

	// If there was an error when stopping at least one of the jobs, return it.
	if hadError {
		return fmt.Errorf("problem stopping one or more job runners - see logs for more info")
	}
	return nil
}
