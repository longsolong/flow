package builtin

import (
	"github.com/longsolong/flow/pkg/workflow/atom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNoopType(t *testing.T) {
	noop := NewNoop("", "")
	assert.Equal(t, "builtin.Noop", noop.Type)
}

func TestNoopRun(t *testing.T) {
	noop := NewNoop("", "")
	err := noop.Create(nil, nil)
	assert.Nil(t, err)
	ret, err := noop.Run(nil)
	assert.Nil(t, err)
	assert.Equal(t, atom.Return{}, ret)
}
