package factory

import (
	"github.com/longsolong/flow/pkg/workflow"
	"github.com/longsolong/flow/pkg/workflow/step/command"
)

func NewEchoCommand(desc string) *command.EchoCommand {
	echo := &command.EchoCommand{}
	echo.Description = desc
	echo.Type = workflow.GenRunnableType(echo)
	return echo
}
