package graph

import (
	"context"
	"github.com/longsolong/flow/pkg/orchestration/job"
	"time"

	"github.com/longsolong/flow/pkg/orchestration/request"
	"github.com/longsolong/flow/pkg/orchestration/standalone/chain"
	"github.com/longsolong/flow/pkg/workflow/atom"
	"github.com/longsolong/flow/pkg/workflow/dag"
)

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

func (p *Plotter) NewNode(ctx context.Context, req *request.Request, step atom.Atom, name string, retry uint, retryWait time.Duration) (*dag.Node, error) {
	if err := step.Create(ctx, req); err != nil {
		return nil, err
	}
	node := dag.NewNode(step, name, retry, retryWait)
	p.DAG.MustAddNode(node)
	p.Chain.AddJob(job.NewJob(step))
	return node, nil
}
