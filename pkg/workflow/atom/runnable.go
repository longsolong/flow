package atom

import (
	"fmt"
	"github.com/longsolong/flow/pkg/workflow"
	"github.com/longsolong/flow/pkg/workflow/state"
	"reflect"
)

// Runnable ...
type Runnable interface {
	Run(ctx workflow.Context) (Return, error)
	Stop(ctx workflow.Context) error
}

// Return ...
type Return struct {
	State state.State // State const
	Exit  int64       // Unix exit code
	Error error       // Go error
}

// GenRunnableType ...
func GenRunnableType(r Runnable, prefix string) string {
	e := reflect.TypeOf(r).Elem()
	return fmt.Sprintf("%s/%s", prefix, e.String())
}

