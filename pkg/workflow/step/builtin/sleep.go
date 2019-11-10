package builtin

import (
	"time"

	"github.com/longsolong/flow/pkg/workflow"
	"github.com/longsolong/flow/pkg/workflow/state"
	"github.com/longsolong/flow/pkg/workflow/step"
)

// Sleep is a Step that sleeps for a given time.
type Sleep struct {
	step.Step

	Duration time.Duration // how long to sleep

	// While running
	stopChan chan struct{}
	stopped  bool
}

// Create ...
func (s *Sleep) Create(ctx workflow.Context) error {
	// TODO parse args from ctx
	durationStr := "1ms"
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return err
	}
	s.Duration = duration
	s.stopChan = make(chan struct{})
	return nil
}

// Run a shell
func (s *Sleep) Run(ctx workflow.Context) (workflow.Return, error) {
	ret := workflow.Return{}

	select {
	case <-time.After(s.Duration):
		ret.State = state.StateSuccess
	case <-s.stopChan:
		ret.State = state.StateCanceled
	}

	return ret, nil
}

// Stop run
func (s *Sleep) Stop(ctx workflow.Context) error {
	if s.stopped {
		return nil
	}
	s.stopped = true

	close(s.stopChan)
	return nil
}