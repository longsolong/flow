package factory

import (
	"testing"

	"github.com/longsolong/flow/pkg/workflow"
	"github.com/stretchr/testify/assert"
)

func TestEchoCommandType(t *testing.T) {
	echo := NewEchoCommand("echo hello")
	assert.Equal(t, "github.com/longsolong/flow/pkg/workflow/step/builtin/command.EchoCommand", echo.Type)
}

func TestEchoCommandRun(t *testing.T) {
	echo := NewEchoCommand("echo hello")
	var ctx workflow.Context
	err := echo.Create(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, echo.Cmd)
	ret, err := echo.Run(ctx)
	assert.Nil(t, err)
	assert.Equal(t, workflow.Return{}, ret)
}
