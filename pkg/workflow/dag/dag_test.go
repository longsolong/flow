package dag

import (
	"github.com/longsolong/flow/pkg/workflow"
	"github.com/longsolong/flow/pkg/workflow/atom"
	"github.com/longsolong/flow/pkg/workflow/step/builtin"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewDAG(t *testing.T) {
	dag := NewDAG("test noop dag", 1)
	noop1 := NewNode(builtin.NewNoop("1", ""), "noop1", 0, time.Duration(0))
	noop2 := NewNode(builtin.NewNoop("2", ""), "noop2", 0, time.Duration(0))

	dag.MustAddNode(noop1)
	dag.MustAddNode(noop2)

	err := noop2.SetUpstream(noop1)
	assert.Nil(t, err)
	err = noop2.SetUpstream(noop1)
	assert.Equal(t, workflow.ErrAlreadyRegisteredUpstream, err)

	upstreams := noop1.Upstream()
	assert.Equal(t, map[atom.AtomID]*Node{}, upstreams)

	_, ok := noop1.Downstream()[noop2.Datum.AtomID()]
	assert.True(t, ok)

	_, ok = noop2.Upstream()[noop1.Datum.AtomID()]
	assert.True(t, ok)

	downstreams := noop2.Downstream()
	assert.Equal(t, map[atom.AtomID]*Node{}, downstreams)
}
