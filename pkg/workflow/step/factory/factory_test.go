package factory

import (
	"testing"

	"github.com/longsolong/flow/pkg/workflow"
	"github.com/longsolong/flow/pkg/workflow/step"
	"github.com/stretchr/testify/assert"
)

func TestEchoCommandType(t *testing.T) {
	echo := NewEchoCommand(step.ID{Name: "echo hello"})
	assert.Equal(t, "builtin/command.EchoCommand", echo.ID.Type)
}

func TestEchoCommandRun(t *testing.T) {
	echo := NewEchoCommand(step.ID{Name: "echo hello"})
	var ctx workflow.Context
	err := echo.Create(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, echo.Cmd)
	ret, err := echo.Run(ctx)
	assert.Nil(t, err)
	assert.Equal(t, workflow.Return{}, ret)
}

func TestSleepRun(t *testing.T) {
	sleep := NewSleep(step.ID{Name: "sleep 1ms"})
	var ctx workflow.Context
	err := sleep.Create(ctx)
	assert.Nil(t, err)
	ret, err := sleep.Run(ctx)
	assert.Nil(t, err)
	assert.Equal(t, workflow.Return{}, ret)
	err = sleep.Stop(ctx)
	assert.Nil(t, err)
}
