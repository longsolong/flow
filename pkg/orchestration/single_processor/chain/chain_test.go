package chain

import (
	"github.com/longsolong/flow/pkg/workflow"
	"github.com/longsolong/flow/pkg/workflow/atom"
	"github.com/longsolong/flow/pkg/workflow/dag"
	"github.com/longsolong/flow/pkg/workflow/step/builtin"
	"time"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewJobChain(t *testing.T) {
	d := dag.NewDAG("test noop chain", 1)

	chain := NewChain(d)
	noop1 := dag.NewNode(
		builtin.NewNoop(atom.ID{ID: "1"}), "noop1", 0, time.Duration(0))
	noop2 := dag.NewNode(
		builtin.NewNoop(atom.ID{ID: "2"}), "noop2", 0, time.Duration(0))

	err := chain.AddNode(noop1)
	assert.Nil(t, err)
	err = chain.AddNode(noop2)
	assert.Nil(t, err)
	err = chain.AddNode(noop2)
	assert.Equal(t, workflow.ErrAlreadyRegisteredNode, err)

	err = noop2.SetUpstream(noop1)
	assert.Nil(t, err)
	err = noop2.SetUpstream(noop1)
	assert.Equal(t, workflow.ErrAlreadyRegisteredUpstream, err)

	upstreams := noop1.Prev
	assert.Equal(t, map[atom.ID]*dag.Node{}, upstreams)

	_, ok := noop1.Next[noop2.Datum.ID()]
	assert.True(t, ok)

	_, ok = noop2.Prev[noop1.Datum.ID()]
	assert.True(t, ok)

	downstreams := noop2.Next
	assert.Equal(t, map[atom.ID]*dag.Node{}, downstreams)
}