package command

import (
	"os/exec"

	"github.com/longsolong/flow/pkg/workflow"
)

// EchoCommand ...
type EchoCommand struct {
	ShellCommand
}

// Create ...
func (echo *EchoCommand) Create(ctx workflow.Context) error {
	echo.Cmd = exec.Command("echo", "hello")
	return nil
}
