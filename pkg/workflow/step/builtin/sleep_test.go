package builtin

import (
	"testing"

	"github.com/longsolong/flow/pkg/workflow/atom"
	"github.com/stretchr/testify/assert"
)

func TestSleepRun(t *testing.T) {
	sleep := NewSleep("", "")
	err := sleep.Create(nil, nil)
	assert.Nil(t, err)
	ret, err := sleep.Run(nil)
	assert.Nil(t, err)
	assert.Equal(t, atom.Return{}, ret)
	err = sleep.Stop(nil)
	assert.Nil(t, err)
}
