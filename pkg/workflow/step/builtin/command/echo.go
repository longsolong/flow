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
	// TODO parse args from ctx
	echo.Cmd = exec.Command("echo", "hello")
	return nil
}
