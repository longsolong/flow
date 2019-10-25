package command

import (
	"github.com/longsolong/flow/pkg/workflow"
	"os/exec"
)

type EchoCommand struct {
	ShellCommand
}

func (echo *EchoCommand) Create(ctx workflow.Context) error {
	echo.Cmd = exec.Command("echo", "hello")
	return nil
}