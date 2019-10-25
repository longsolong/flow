package command

import (
	"os/exec"

	"github.com/longsolong/flow/pkg/workflow"
	"github.com/longsolong/flow/pkg/workflow/state"
	"github.com/longsolong/flow/pkg/workflow/step"
)

// ShellCommand is a step that runs a single shell command with arguments.
type ShellCommand struct {
	step.Step

	Cmd *exec.Cmd
}

// Run a shell
func (s *ShellCommand) Run(ctx workflow.Context) (workflow.Return, error) {
	// Run the cmd and wait for it to return
	exit := int64(0)
	err := s.Cmd.Run()
	ret := workflow.Return{
		Exit:  exit,
		Error: err,
	}
	if err != nil {
		ret.Exit = 1
		ret.State = state.STATE_FAIL
	} else {
		ret.State = state.STATE_SUCCESS
	}

	return ret, nil
}

// Stop run
func (s *ShellCommand) Stop(ctx workflow.Context) error {
	return nil
}
