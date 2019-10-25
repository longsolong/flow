package workflow

import (
	"fmt"
	"reflect"
)

// Runnable ...
type Runnable interface {
	Run(ctx Context) (Return, error)
	Stop(ctx Context) error
}

// Return ...
type Return struct {
	State byte  // STATE_ const
	Exit  int64 // Unix exit code
	Error error // Go error
}

// GenRunnableType ...
func GenRunnableType(r Runnable) string {
	e := reflect.TypeOf(r).Elem()
	return fmt.Sprintf("%s.%s", e.PkgPath(), e.Name())
}
