package step

import (
	"fmt"
	"github.com/longsolong/flow/pkg/workflow"
	"reflect"
)

type Step interface {
	Create(ctx workflow.Context) error

	Run(ctx workflow.Context) (Return, error)

	Stop() error

	Status() string
}

type Return struct {
	State  byte   // STATE_ const
	Exit   int64  // Unix exit code
	Error  error  // Go error
}


func GenType(s Step) string {
	empty := reflect.TypeOf(s).Elem()
	return fmt.Sprintf("%s.%s", empty.PkgPath(), empty.Name())
}

type Atom struct {
	Type string
	status string
}

func (atom *Atom) Create(ctx workflow.Context) error {
	args, err := atom.NewArgs(ctx)
	if err != nil {
		return err
	}
	return atom.New(args...)
}

func (atom *Atom) Run(ctx workflow.Context) error {
	return nil
}

func (atom *Atom) NewArgs(ctx workflow.Context) ([]interface{}, error) {
	return nil, nil
}

func (atom *Atom) New(arg ...interface{}) error {
	return nil
}

func (atom *Atom) Stop() error {
	return nil
}

func (atom *Atom) Status() string {
	return atom.status
}

func (atom *Atom) SetStatus(msg string) {
	atom.status = msg
}