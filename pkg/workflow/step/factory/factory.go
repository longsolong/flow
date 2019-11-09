package factory

import (
	"github.com/longsolong/flow/pkg/workflow"
	"github.com/longsolong/flow/pkg/workflow/step/builtin"
	"github.com/longsolong/flow/pkg/workflow/step/builtin/command"
)

// NewEchoCommand ...
func NewEchoCommand(name string) *command.EchoCommand {
	echo := &command.EchoCommand{}
	echo.ID.Name = name
	echo.ID.Type = workflow.GenRunnableType(echo, "builtin")
	return echo
}

// NewSleep ...
func NewSleep(name string) *builtin.Sleep {
	s := &builtin.Sleep{}
	s.ID.Name = name
	s.ID.Type = workflow.GenRunnableType(s, "builtin")
	return s
}
