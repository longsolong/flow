package graph

import (
	"github.com/longsolong/flow/pkg/orchestration/request"
	"github.com/longsolong/flow/pkg/orchestration/single_processor/chain"
	"github.com/longsolong/flow/pkg/workflow/dag"
)

// Grapher ...
type Grapher struct {
	Req *request.Request

	*dag.DAG
	*chain.Chain
	GraphPlotter
}

// NewGrapher A new Grapher should be made for every request.
func NewGrapher(
	req *request.Request, d *dag.DAG, c *chain.Chain, p GraphPlotter) *Grapher {

	g := &Grapher{
		Req: req,
		DAG: d,
		Chain: c,
		GraphPlotter: p,
	}
	return g
}
