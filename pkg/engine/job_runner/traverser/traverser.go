package traverser

import "context"

// A Traverser provides the ability to run a job traverser while respecting the
// dependencies between the jobs.
type Traverser interface {
	// Run traverses a job traverser and runs all of the jobs in it. It starts by
	// running the first job in the traverser, and then, if the job completed,
	// successfully, running its adjacent jobs. This process continues until there
	// are no more jobs to run, or until the Stop method is called on the traverser.
	Run(ctx context.Context)

	// Stop makes a traverser stop traversing its job traverser. It also sends a stop
	// signal to all of the jobs that a traverser is running.
	//
	// It returns an error if it fails to stop all running jobs.
	Stop(ctx context.Context) error
}
