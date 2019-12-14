package command

import (
	"context"
	"github.com/longsolong/flow/pkg/orchestration/request"
	"github.com/longsolong/flow/pkg/workflow/atom"
	"os/exec"
)

// EchoCommand ...
type EchoCommand struct {
	ShellCommand
}

// NewEchoCommand ...
func NewEchoCommand(id, expansionDigest string) *EchoCommand {
	echo := &EchoCommand{}
	echo.SetID(atom.ID{
		ID:              id,
		ExpansionDigest: expansionDigest,
		Type:            atom.GenRunnableType(echo, "builtin"),
	})
	return echo
}

// Create ...
func (echo *EchoCommand) Create(ctx context.Context, req *request.Request) error {
	// TODO parse args from ctx
	echo.Cmd = exec.Command("echo", "hello")
	return nil
}
