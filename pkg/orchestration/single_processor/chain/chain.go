package chain

import (
	"github.com/longsolong/flow/pkg/orchestration/job_runner/chain"
	"github.com/longsolong/flow/pkg/workflow/dag"
)

// Chain represents a job chain and some meta information about it.
type Chain struct {
	*chain.JobChain
}

// NewChain ...
func NewChain(d *dag.DAG) *Chain {
	return &Chain{
		JobChain: chain.NewJobChain(d),
	}
}
