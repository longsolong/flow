package command

import (
	"github.com/google/uuid"
	"github.com/longsolong/flow/pkg/workflow/step"
)

type EchoCommand struct {
	ShellCommand
}

func NewEchoCommand(uuid uuid.UUID) *EchoCommand {
	id := step.Id{Uuid: uuid, Type: step.GenType((*EchoCommand)(nil))}
	echo := &EchoCommand{}
	echo.SetId(id)
	return echo
}
