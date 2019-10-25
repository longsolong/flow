package factory

import (
	"github.com/longsolong/flow/pkg/workflow"
	"github.com/longsolong/flow/pkg/workflow/step/builtin/command"
)

// NewEchoCommand ...
func NewEchoCommand(desc string) *command.EchoCommand {
	echo := &command.EchoCommand{}
	echo.Description = desc
	echo.Type = workflow.GenRunnableType(echo)
	return echo
}
