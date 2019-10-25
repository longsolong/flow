package command

import (
	"fmt"
	"github.com/longsolong/flow/pkg/workflow"
	"github.com/longsolong/flow/pkg/workflow/state"
	"github.com/longsolong/flow/pkg/workflow/step"
	"os/exec"
)

// ShellCommand is a step that runs a single shell command with arguments.
type ShellCommand struct {
	step.Step

	Cmd *exec.Cmd

	status string
}

func (s *ShellCommand) Run(ctx workflow.Context) (workflow.Return, error) {
	// Set status before and after
	s.SetStatus(fmt.Sprintf("runnning %v", s.Cmd.Args))
	defer s.SetStatus(fmt.Sprintf("done running %v", s.Cmd.Args))

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

func (s *ShellCommand) Stop() error {
	return nil
}

func (s *ShellCommand) Status() string {
	return s.status
}

func (s *ShellCommand) SetStatus(msg string) {
	s.status = msg
}