package graph

import (
	"context"

	"github.com/longsolong/flow/pkg/orchestration/request"
	"github.com/longsolong/flow/pkg/orchestration/single_processor/chain"
	"github.com/longsolong/flow/pkg/workflow/dag"
)

// GraphPlotter ...
type GraphPlotter interface {
	Begin(ctx context.Context, req *request.Request) error
	Grow(ctx context.Context)
	Done() <-chan struct{}
}

// Plotter ...
type Plotter struct {
	*dag.DAG
	*chain.Chain

	done chan struct{}
}

func NewPlotter(name string, version int) (p Plotter) {
	p.done = make(chan struct{})
	p.DAG = dag.NewDAG(name, version)
	p.Chain = chain.NewChain(p.DAG)
	return p
}

func (p *Plotter) Done() <-chan struct{} {
	return p.done
}

func (p *Plotter) Close() {
	close(p.done)
}
