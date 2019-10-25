package command

import (
	"github.com/google/uuid"
	"github.com/longsolong/flow/pkg/workflow/step"
)

type EchoCommand struct {
	ShellCommand
}

func NewEchoCommand(uuid uuid.UUID) *EchoCommand {
	echo := &EchoCommand{}
	echo.Type = step.GenType((*EchoCommand)(nil))
	return echo
}
