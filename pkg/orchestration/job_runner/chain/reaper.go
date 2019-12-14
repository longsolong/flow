package chain

// JobReaper handles jobs and chains that have finished running.
type JobReaper interface {
	// Run reaps done jobs from doneJobChan, saving their states and enqueueing
	// any jobs that should be run to runJobChan. When there are no more jobs to
	// reap, Run finalizes the chain and returns.
	Run()

	// Stop stops the JobReaper from reaping any more jobs. It blocks until
	// Run() returns and the reaper can be safely switched out for another
	// implementation.
	Stop()
}
