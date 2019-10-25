package workflow

import (
	"fmt"
	"reflect"
)

type Runnable interface {
	Run(ctx Context) (Return, error)

	Stop() error

	Status() string
}

type Return struct {
	State byte  // STATE_ const
	Exit  int64 // Unix exit code
	Error error // Go error
}

func GenRunnableType(r Runnable) string {
	e := reflect.TypeOf(r).Elem()
	return fmt.Sprintf("%s.%s", e.PkgPath(), e.Name())
}

