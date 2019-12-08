package command

import (
	"github.com/longsolong/flow/pkg/workflow"
	"github.com/longsolong/flow/pkg/workflow/atom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEchoCommandType(t *testing.T) {
	echo := NewEchoCommand("", "")
	assert.Equal(t, "builtin/command.EchoCommand", echo.ID().Type)
}


func TestEchoCommandRun(t *testing.T) {
	echo := NewEchoCommand("", "")
	var ctx workflow.Context
	err := echo.Create(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, echo.Cmd)
	ret, err := echo.Run(ctx)
	assert.Nil(t, err)
	assert.Equal(t, atom.Return{}, ret)
}

