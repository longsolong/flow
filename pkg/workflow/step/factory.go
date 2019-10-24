package step

import (
	"fmt"
	"github.com/longsolong/flow/pkg/workflow"
	"github.com/longsolong/flow/pkg/workflow/state"

	"os/exec"
)

// ShellCommand is a step that runs a single shell command with arguments.
type ShellCommand struct {
	*Atom

	Cmd *exec.Cmd
}

func (s *ShellCommand) New(name string, arg ...string) error {
	s.Cmd = exec.Command(name, arg...)
	return nil
}

func (s *ShellCommand) Run(ctx workflow.Context) (Return, error) {
	// Set status before and after
	s.setStatus(fmt.Sprintf("runnning %v",  s.Cmd.Args))
	defer s.setStatus(fmt.Sprintf("done running %v",  s.Cmd.Args))

	// Run the cmd and wait for it to return
	exit := int64(0)
	err := s.Cmd.Run()
	ret := Return{
		Exit:   exit,
		Error:  err,
	}
	if err != nil {
		ret.Exit = 1
		ret.State = state.STATE_FAIL
	} else {
		ret.State = state.STATE_SUCCESS
	}

	return ret, nil
}