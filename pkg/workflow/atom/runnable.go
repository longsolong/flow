package atom

import (
	"context"
	"github.com/longsolong/flow/pkg/workflow/state"
)

// Runnable ...
type Runnable interface {
	Run(ctx context.Context) (Return, error)
	Stop(ctx context.Context) error
}

// Return ...
type Return struct {
	State state.State // State const
	Exit  int64       // Unix exit code
	Error error       // Go error
}