package builtin

import (
	"github.com/longsolong/flow/pkg/workflow"
	"github.com/longsolong/flow/pkg/workflow/atom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSleepRun(t *testing.T) {
	sleep := NewSleep("", "")
	var ctx workflow.Context
	err := sleep.Create(ctx)
	assert.Nil(t, err)
	ret, err := sleep.Run(ctx)
	assert.Nil(t, err)
	assert.Equal(t, atom.Return{}, ret)
	err = sleep.Stop(ctx)
	assert.Nil(t, err)
}

