package builtin

import (
	"context"
	"github.com/longsolong/flow/pkg/orchestration/request"

	"github.com/longsolong/flow/pkg/workflow/atom"
	"github.com/longsolong/flow/pkg/workflow/step"
)

// Noop ...
type Noop struct {
	step.Step
}

// NewNoop ...
func NewNoop(id atom.ID) *Noop {
	n := &Noop{}
	id.Type = atom.GenRunnableType(n, "builtin")
	n.SetID(id)
	return n
}

// Create ...
func (s *Noop) Create(ctx context.Context, req *request.Request) error {
	return nil
}

// Run ...
func (s *Noop) Run(ctx context.Context) (atom.Return, error) {
	ret := atom.Return{}
	return ret, nil
}

// Stop run
func (s *Noop) Stop(ctx context.Context) error {
	return nil
}
