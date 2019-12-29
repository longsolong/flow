package builtin

import (
	"context"
	"github.com/longsolong/flow/pkg/orchestration/request"
	"time"

	"github.com/longsolong/flow/pkg/workflow/atom"
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

// NewSleep ...
func NewSleep(id, expansionDigest string) *Sleep {
	s := &Sleep{}
	s.SetID(atom.ID{
		ID:              id,
		ExpansionDigest: expansionDigest,
		Type:            atom.GenRunnableType(s, "builtin"),
	})
	return s
}

// Create ...
func (s *Sleep) Create(ctx context.Context, req *request.Request) error {
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

// Run ...
func (s *Sleep) Run(ctx context.Context) (atom.Return, error) {
	ret := atom.Return{}

	select {
	case <-time.After(s.Duration):
		ret.State = state.StateSuccess
	case <-s.stopChan:
		ret.State = state.StateCanceled
	}

	return ret, nil
}

// Stop run
func (s *Sleep) Stop(ctx context.Context) error {
	if s.stopped {
		return nil
	}
	s.stopped = true

	close(s.stopChan)
	return nil
}
