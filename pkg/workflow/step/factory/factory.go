package factory

import (
	"github.com/longsolong/flow/pkg/workflow"
	"github.com/longsolong/flow/pkg/workflow/step"
	"github.com/longsolong/flow/pkg/workflow/step/builtin"
	"github.com/longsolong/flow/pkg/workflow/step/builtin/command"
)

// NewEchoCommand ...
func NewEchoCommand(id step.ID) *command.EchoCommand {
	echo := &command.EchoCommand{}
	echo.ID = id
	echo.ID.Type = workflow.GenRunnableType(echo, "builtin")
	return echo
}

// NewSleep ...
func NewSleep(id step.ID) *builtin.Sleep {
	s := &builtin.Sleep{}
	s.ID = id
	s.ID.Type = workflow.GenRunnableType(s, "builtin")
	return s
}

// NewNoop ...
func NewNoop(id step.ID) *builtin.Noop {
	n := &builtin.Noop{}
	n.ID = id
	n.ID.Type = workflow.GenRunnableType(n, "builtin")
	return n
}
