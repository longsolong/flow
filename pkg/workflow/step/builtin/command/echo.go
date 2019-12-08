package command

import (
	"github.com/longsolong/flow/pkg/workflow/atom"
	"os/exec"
	"github.com/longsolong/flow/pkg/workflow"
)

// EchoCommand ...
type EchoCommand struct {
	ShellCommand
}

// NewEchoCommand ...
func NewEchoCommand(id, expansionDigest string) *EchoCommand {
	echo := &EchoCommand{}
	echo.SetID(atom.ID{
		ID: id,
		ExpansionDigest: expansionDigest,
		Type: atom.GenRunnableType(echo, "builtin"),
	})
	return echo
}

// Create ...
func (echo *EchoCommand) Create(ctx workflow.Context) error {
	// TODO parse args from ctx
	echo.Cmd = exec.Command("echo", "hello")
	return nil
}
