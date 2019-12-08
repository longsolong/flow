package command

import (
	"github.com/longsolong/flow/pkg/workflow/atom"
	"os/exec"

	"github.com/longsolong/flow/pkg/workflow"
	"github.com/longsolong/flow/pkg/workflow/state"
	"github.com/longsolong/flow/pkg/workflow/step"
)

// ShellCommand is a Step that runs a single shell command with arguments.
type ShellCommand struct {
	step.Step

	Cmd *exec.Cmd
}

// Run a shell
func (s *ShellCommand) Run(ctx workflow.Context) (atom.Return, error) {
	// Run the cmd and wait for it to return
	err := s.Cmd.Run()
	ret := atom.Return{
		Error: err,
	}
	if err != nil {
		ret.Exit = 1
		ret.State = state.StateFail
	} else {
		ret.State = state.StateSuccess
	}

	return ret, nil
}

// Stop run
func (s *ShellCommand) Stop(ctx workflow.Context) error {
	return nil
}
