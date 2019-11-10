package dag

import (
	"github.com/longsolong/flow/pkg/workflow"
	"github.com/longsolong/flow/pkg/workflow/step"
	"github.com/longsolong/flow/pkg/workflow/step/factory"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStepChain(t *testing.T) {
	chain := NewDAG()
	noop1 := factory.NewNoop(step.ID{Name: "noop1", ID: "1"})
	noop2 := factory.NewNoop(step.ID{Name: "noop2", ID: "2"})

	err := chain.AddStep(noop1.ID, noop1)
	assert.Nil(t, err)
	err = chain.AddStep(noop2.ID, noop2)
	assert.Nil(t, err)
	err = chain.AddStep(noop2.ID, noop2)
	assert.Equal(t, workflow.ErrAlreadyRegisteredStep, err)

	err = chain.SetUpstream(noop2.ID, noop1.ID)
	assert.Nil(t, err)
	err = chain.SetUpstream(noop2.ID, noop1.ID)
	assert.Equal(t, workflow.ErrAlreadyRegisteredUpstream, err)

	upstreams := chain.Upstreams(noop1.ID)
	assert.Equal(t, ([]step.Atom)(nil), upstreams)

	upstreams = chain.Upstreams(noop2.ID)
	assert.Equal(t, []step.Atom{noop1}, upstreams)

}
