package dag

import (
	"github.com/longsolong/flow/pkg/workflow"
	"github.com/longsolong/flow/pkg/workflow/step"
)

// DAG represents a directed acyclic graph.
type DAG struct {
	Steps           map[step.ID]step.Atom
	UpstreamStepIDs map[step.ID]map[step.ID]bool
}

// NewDAG ...
func NewDAG() *DAG {
	return &DAG{
		Steps:           make(map[step.ID]step.Atom),
		UpstreamStepIDs: make(map[step.ID]map[step.ID]bool),
	}
}

// AddStep ...
func (g *DAG) AddStep(stepID step.ID, atom step.Atom) error {
	if _, ok := g.Steps[stepID]; ok {
		return workflow.ErrAlreadyRegisteredStep
	}
	g.Steps[stepID] = atom
	return nil
}

// SetUpstream ...
func (g *DAG) SetUpstream(currentStepID, upstreamStepID step.ID) error {
	upstreams, ok := g.UpstreamStepIDs[currentStepID]
	if ok {
		if _, ok2 := upstreams[upstreamStepID]; ok2 {
			return workflow.ErrAlreadyRegisteredUpstream
		}
	} else {
		g.UpstreamStepIDs[currentStepID] = make(map[step.ID]bool)
	}
	g.UpstreamStepIDs[currentStepID][upstreamStepID] = true
	return nil
}

// Upstreams ...
func (g *DAG) Upstreams(stepID step.ID) (upstreams []step.Atom) {
	if upstreamStepIDs, ok := g.UpstreamStepIDs[stepID]; ok {
		for stepID := range upstreamStepIDs {
			upstreams = append(upstreams, g.Steps[stepID])
		}
	}
	return
}
