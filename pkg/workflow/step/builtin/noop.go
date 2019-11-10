package builtin

import (
	"github.com/longsolong/flow/pkg/workflow"
	"github.com/longsolong/flow/pkg/workflow/step"
)

// Noop ...
type Noop struct {
	step.Step
}

// Create ...
func (s *Noop) Create(ctx workflow.Context) error {
	return nil
}

// Run ...
func (s *Noop) Run(ctx workflow.Context) (workflow.Return, error) {
	ret := workflow.Return{}
	return ret, nil
}

// Stop run
func (s *Noop) Stop(ctx workflow.Context) error {
	return nil
}
