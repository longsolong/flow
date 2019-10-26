package factory

import (
	"github.com/longsolong/flow/pkg/workflow"
	"github.com/longsolong/flow/pkg/workflow/step/builtin"
	"github.com/longsolong/flow/pkg/workflow/step/builtin/command"
)

// NewEchoCommand ...
func NewEchoCommand(desc string) *command.EchoCommand {
	echo := &command.EchoCommand{}
	echo.Description = desc
	echo.Type = workflow.GenRunnableType(echo, "builtin")
	return echo
}

// NewSleep ...
func NewSleep(desc string) *builtin.Sleep {
	s := &builtin.Sleep{}
	s.Description = desc
	s.Type = workflow.GenRunnableType(s, "builtin")
	return s
}
