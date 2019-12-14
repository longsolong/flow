package graph

import (
	"github.com/longsolong/flow/pkg/orchestration/request"
	"github.com/longsolong/flow/pkg/orchestration/single_processor/chain"
	"github.com/longsolong/flow/pkg/workflow/dag"
)

// Plotter ...
type Plotter interface {
	Begin(name string, version int, req *request.Request) (*dag.DAG, *chain.Chain, error)
	Grow(d *dag.DAG, c *chain.Chain, req *request.Request) error
}