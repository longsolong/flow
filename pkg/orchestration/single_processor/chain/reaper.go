package chain

import (
	"github.com/longsolong/flow/pkg/infra"
	"github.com/longsolong/flow/pkg/orchestration/job"
	"github.com/longsolong/flow/pkg/workflow/atom"
	"github.com/longsolong/flow/pkg/workflow/state"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

type reaper struct {
	chain  *Chain
	logger *infra.Logger

	stopMux  *sync.Mutex
	stopped  bool
	stopChan chan struct{}

	doneChan    chan struct{}
	doneJobChan chan job.Job
}

// RunningChainReaper ...
type RunningChainReaper struct {
	reaper
	runJobChan chan job.Job // enqueue next jobs to run here
}

// Run reaps jobs when they finish running. For each job reaped, if...
// - chain is done: save final state .
// - job failed:    retry sequence if possible.
// - job completed: prepared subsequent jobs and enqueue if runnable.
func (r *RunningChainReaper) Run() {
	defer close(r.doneChan)

	// If the chain is already done, skip straight to finalizing.
	done, complete := r.chain.IsDoneRunning()
	if done {
		r.Finalize(complete)
		return
	}

REAPER:
	for {
		select {
		case j := <-r.doneJobChan:
			r.Reap(&j)
			done, complete = r.chain.IsDoneRunning()
			if done {
				break REAPER
			}
		case <-r.stopChan:
			// Don't Finalize the chain when stopping
			return
		}
	}

	r.Finalize(complete)
}

// Stop stops the reaper from reaping any more jobs. It blocks until the reaper
// is stopped (will reap no more jobs and Run will return).
func (r *RunningChainReaper) Stop() {
	r.stopMux.Lock()
	defer r.stopMux.Unlock()
	if r.stopped {
		return
	}
	r.stopped = true

	close(r.stopChan)
	<-r.doneChan
	return
}

// Reap takes a job that just finished running, saves its final state, and prepares
// to continue running the chain (or recognizes that the chain is done running).
//
// If chain is done: save final state + stop running more jobs.
// If job failed:    retry sequence if possible.
// If job completed: prepared subsequent jobs and enqueue if runnable.
func (r *RunningChainReaper) Reap(job *job.Job) {
	fields := []zapcore.Field{
		zap.String("job_id", job.ID().String()),
		zap.String("sequence_id", r.chain.DAG.Vertices[job.ID()].SequenceID.String()),
		zap.Int("sequence_try", int(r.chain.DAG.Vertices[job.ID()].SequenceRetry)),
	}
	logger := r.logger.Log

	// Set the final state of the job in the chain.
	r.chain.SetJobState(job.ID(), job.State)

	if _, ok := state.JobCompleteState[job.State]; ok {
		for _, nextJob := range r.chain.NextJobs(job.ID()) {
			nextFields := append([]zapcore.Field(nil), fields...)

			if !r.chain.IsRunnable(nextJob.ID()) {
				logger.Info("next job not runnable", nextFields...)
				continue
			}
			logger.Info("enqueueing next job", nextFields...)
			r.runJobChan <- *nextJob
		}
	} else {
		// Job was NOT successful. The job.Runner already did job retries.
		// Retry sequence if possible.
		if !r.chain.CanRetrySequence(job.ID()) {
			logger.Warn("job failed, no sequence tries left", fields...)
			return
		}
		logger.Warn("job failed, retrying sequence", fields...)
		sequenceStartJob := r.prepareSequenceRetry(job)
		r.runJobChan <- *sequenceStartJob // re-enqueue first job in sequence
	}
}

// Finalize determines the final state of the chain
func (r *RunningChainReaper) Finalize(complete bool) {
}

// prepareSequenceRetry prepares a sequence to retry. The caller should check
// r.chain.CanRetrySequence first; this func does not check the seq retry limit
// or increment seq try count (that's done in traverser.runJobs when the seq
// start job runs).
func (r *reaper) prepareSequenceRetry(failedJob *job.Job) *job.Job {
	sequenceStartJob := r.chain.SequenceStartJob(failedJob.ID())

	fields := []zapcore.Field{
		zap.String("sequence_id", sequenceStartJob.ID().String()),
	}
	logger := r.logger.Log

	logger.Info("preparing sequence retry", fields...)

	// sequenceJobsToRetry is a list containing the failed job and all previously
	// completed jobs in the sequence. For example, if job C of A -> B -> C -> D
	// fails, then A and B are the previously completed jobs and C is the failed
	// job. So, jobs A, B, and C will be added to sequenceJobsToRetry. D will not be
	// added because it was never run.
	sequenceJobsToRetry := r.sequenceJobsCompleted(sequenceStartJob)

	haveFailedJob := false
	for _, j := range sequenceJobsToRetry {
		if j.ID() == failedJob.ID() {
			haveFailedJob = true
			break
		}
	}
	if !haveFailedJob {
		sequenceJobsToRetry = append(sequenceJobsToRetry, failedJob)
	}

	// Roll back completed sequence jobs
	finishedJobs := 0
	for _, j := range sequenceJobsToRetry {
		jobState := r.chain.JobState(j.ID())
		if _, ok := state.JobCompleteState[jobState]; ok {
			finishedJobs++
		}

		// Roll back job state to pending so it's runnable again
		r.chain.SetJobState(j.ID(), state.StateUpForRetry)
	}

	// Running reaper will re-enqueue/re-run seq from this seq start job.
	// Suspend reaper will not, leaving seq in runnable state for when chain is resumed.
	return sequenceStartJob
}

// sequenceJobsCompleted does a BFS to find all jobs in the sequence that have
// completed. You can read how BFS works here:
// https://en.wikipedia.org/wiki/Breadth-first_search.
func (r *reaper) sequenceJobsCompleted(sequenceStartJob *job.Job) []*job.Job {
	toVisit := map[atom.ID]*job.Job{} // job id -> job to visit
	visited := map[atom.ID]*job.Job{} // job id -> job visited

	// Process sequenceStartJob
	for _, pJob := range r.chain.NextJobs(sequenceStartJob.ID()) {
		toVisit[pJob.ID()] = pJob
	}
	visited[sequenceStartJob.ID()] = sequenceStartJob

PROCESS_TO_VISIT_LIST:
	for len(toVisit) > 0 {

	PROCESS_CURRENT_JOB:
		for currentJobID, currentJob := range toVisit {

		PROCESS_NEXT_JOBS:
			for _, nextJob := range r.chain.NextJobs(currentJobID) {
				// Don't add failed or pending jobs to toVisit list
				// For example, if job C of A -> B -> C -> D fails, then do not add C
				// or D to toVisit list. Because we have single sequence retries,
				// stopping at the failed job ensures we do not add jobs not in the
				// sequence to the toVisit list.
				if _, ok := state.JobCompleteState[nextJob.State]; !ok {
					continue PROCESS_NEXT_JOBS
				}

				// Make sure we don't visit a job multiple times. We can see a job
				// multiple times if it is a "fan in" node.
				if _, seen := visited[nextJob.ID()]; !seen {
					toVisit[nextJob.ID()] = nextJob
				}
			}

			// Since we have processed all of the next jobs for this current job, we
			// are done visiting the current job and can delete it from the toVisit
			// list and add it to the visited list.
			delete(toVisit, currentJobID)
			visited[currentJobID] = currentJob

			continue PROCESS_CURRENT_JOB
		}

		continue PROCESS_TO_VISIT_LIST
	}

	completedJobs := make([]*job.Job, 0, len(visited))
	for _, j := range visited {
		completedJobs = append(completedJobs, j)
	}

	return completedJobs
}
