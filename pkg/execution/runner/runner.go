// Copyright 2017-2019, Square, Inc.

package runner

import (
	"context"
	"fmt"
	"github.com/longsolong/flow/pkg/infra"
	"github.com/longsolong/flow/pkg/orchestration/job"
	"github.com/longsolong/flow/pkg/orchestration/request"
	"github.com/longsolong/flow/pkg/workflow/atom"
	flowcontext "github.com/longsolong/flow/pkg/workflow/context"
	"github.com/longsolong/flow/pkg/workflow/state"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"runtime/debug"
	"sync"
	"time"
)

// Return ...
type Return struct {
	AtomReturn atom.Return // Final atom.Return Determines if/how chain continues running.
	Tries      uint        // Number of tries this run, not including any previous tries
}

// A Runner runs and manages one job in a job chain. The job must implement the
// job.Job interface.
type Runner interface {
	// Run runs the job, blocking until it has completed or when Stop is called.
	// If the job fails, Run will retry it as many times as the job is configured
	// to be retried. When the job successfully completes, or reaches the maximum number
	// of retry attempts, Run returns the final state of the job.
	Run(ctx context.Context) Return

	// Stop stops the job if it's running. The job is responsible for stopping
	// quickly because Stop blocks while waiting for the job to stop.
	Stop(ctx context.Context) error
}

// A runner represents all information needed to run a job.
type runner struct {
	jobName string
	realJob job.Job // the actual job interface to run
	req     *request.Request

	totalTries uint // try count all seq tries
	maxTries   uint // max tries per seq try
	retryWait  time.Duration
	stopChan   chan struct{}
	mux        *sync.Mutex

	logger    *infra.Logger
	startTime time.Time
	sleeping  bool
}

// NewRunner ...
func NewRunner(realJob job.Job, req *request.Request, totalTries uint, name string, retry uint, retryWait time.Duration, logger *infra.Logger) Runner {
	return &runner{
		jobName: name,
		realJob: realJob,
		req:     req,

		totalTries: 1 + totalTries, // this run + past totalTries (on resume/retry)
		maxTries:   1 + retry,      // + 1 because we always run once
		retryWait:  retryWait,
		stopChan:   make(chan struct{}),
		mux:        &sync.Mutex{},

		logger:    logger,
		startTime: time.Now().UTC(),
	}
}

func (r *runner) Run(ctx context.Context) Return {
	fields := []zapcore.Field{
		zap.String("request_id", r.req.RequestUUID.String()),
		zap.String("job_id", r.realJob.AtomID().String()),
	}
	logger := r.logger.Log

	// The chain.traverser that's calling us only cares about the final state
	// of the job. If maxTries > 1, the intermediate states are only logged if
	// the run fails.
	tries := uint(1) // number of tries this run
	finalAtomReturn := atom.Return{State: state.StateUnknown}
TRY_LOOP:
	for tries <= r.maxTries {
		tryFields := append([]zapcore.Field(nil), fields...)
		tryFields = append(tryFields, []zapcore.Field{
			zap.Uint("try", r.totalTries),
			zap.Uint("tries", tries),
			zap.Uint("max_tries", r.maxTries),
		}...)

		// Can be stopped before we've started.
		if r.stopped() {
			logger.Info("job stopped before start", tryFields...)
			break TRY_LOOP
		}

		// Run the job. Use a separate method so we can easily recover from a panic
		// in job.Run.
		logger.Info("job start", tryFields...)
		ctx = context.WithValue(ctx, flowcontext.FlowContextKey("totalTries"), r.totalTries)
		ctx = context.WithValue(ctx, flowcontext.FlowContextKey("maxTries"), r.maxTries)
		startedAt, finishedAt, jobRet, runErr := r.runJob(ctx)
		runtime := time.Duration(finishedAt-startedAt) * time.Nanosecond
		runFields := append([]zapcore.Field(nil), fields...)
		runFields = append(runFields, []zapcore.Field{
			zap.Duration("runtime", runtime),
			zap.String("state", state.StateText[jobRet.State]),
			zap.Int64("exit", jobRet.Exit),
			zap.NamedError("err", jobRet.Error),
			zap.NamedError("run_err", runErr),
		}...)
		logger.Info("job return", runFields...)

		// Set final job state to this job state
		finalAtomReturn.State = jobRet.State

		// Break try loop on success or stop
		if jobRet.State == state.StateSuccess || jobRet.State == state.StateStopped {
			break TRY_LOOP
		}

		// Job failed, wait and retry?
		retryFields := append([]zapcore.Field(nil), fields...)
		retryFields = append(retryFields, []zapcore.Field{
			zap.Uint("try_left", r.maxTries-tries),
		}...)

		logger.Warn("job failed", retryFields...)

		// If last try, break retry loop, don't wait
		if tries == r.maxTries {
			break TRY_LOOP
		}

		// Wait between retries. Can be stopped while waiting which is why we
		// need to increment tries first. At this point, we're effectively on
		// the next try.
		r.mux.Lock()
		tries++
		r.totalTries++
		r.sleeping = true
		r.mux.Unlock()
		select {
		case <-time.After(r.retryWait):
			// Job failed, wait and retry?
			r.mux.Lock()
			r.sleeping = false
			r.mux.Unlock()
		case <-r.stopChan:
			waitRetryFields := append([]zapcore.Field(nil), fields...)
			waitRetryFields = append(waitRetryFields, []zapcore.Field{
				zap.Uint("total_tries", r.totalTries),
			}...)
			logger.Info("job stopped while waiting to run try", waitRetryFields...)
			break TRY_LOOP
		}
	}

	return Return{
		AtomReturn: finalAtomReturn,
		Tries:      tries,
	}
}

// Actually run the job.
func (r *runner) runJob(ctx context.Context) (startedAt, finishedAt int64, ret atom.Return, err error) {
	defer func() {
		// Recover from a panic inside Job.Run()
		if panicErr := recover(); panicErr != nil {
			// Set named return values. startedAt will already be set before
			// the panic.
			finishedAt = time.Now().UnixNano()
			ret = atom.Return{
				State: state.StateException,
				Exit:  1,
			}
			// The returned error will be used in the job log entry.
			err = fmt.Errorf("panic from job.Run: %s", panicErr)
			debug.PrintStack()
		}
	}()

	// Run the job. Run is a blocking operation that could take a long
	// time. Run will return when a job finishes running (either by
	// its own accord or by being forced to finish when Stop is called).
	startedAt = time.Now().UnixNano()
	ret, err = r.realJob.Atom.(atom.Runnable).Run(ctx)
	finishedAt = time.Now().UnixNano()

	return
}

func (r *runner) Stop(ctx context.Context) error {
	fields := []zapcore.Field{
		zap.String("request_id", r.req.RequestUUID.String()),
		zap.String("job_id", r.realJob.AtomID().String()),
	}
	logger := r.logger.Log

	r.mux.Lock() // LOCK

	// Return if stop was already called.
	select {
	case <-r.stopChan:
		r.mux.Unlock() // UNLOCK
		return nil
	default:
	}

	close(r.stopChan)

	r.mux.Unlock() // UNLOCK

	logger.Info("stopping the job", fields...)
	return r.realJob.Atom.(atom.Runnable).Stop(ctx) // this is a blocking operation that should return quickly
}

func (r *runner) stopped() bool {
	select {
	case <-r.stopChan:
		return true
	default:
		return false
	}
}
