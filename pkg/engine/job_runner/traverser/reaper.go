package traverser

import "context"

// JobReaper handles jobs and chains that have finished running.
type JobReaper interface {
	// Run reaps done jobs from doneJobChan, saving their states and enqueueing
	// any jobs that should be run to runJobChan. When there are no more jobs to
	// reap, Run finalizes the traverser and returns.
	Run(ctx context.Context)

	// Stop stops the JobReaper from reaping any more jobs. It blocks until
	// Run() returns.
	Stop(ctx context.Context)
}
