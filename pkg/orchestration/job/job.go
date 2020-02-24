package job

import (
	"encoding/json"
	"github.com/longsolong/flow/pkg/workflow/atom"
	"github.com/longsolong/flow/pkg/workflow/state"
)

// Job ...
type Job struct {
	atom.Atom             // step
	State     state.State // State const
}

// NewJob ...
func NewJob(step atom.Atom) *Job {
	return &Job{
		Atom:        step,
		State:       state.StateUnknown,
	}
}

func (j Job) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Atom": j.Atom,
		"State": j.State,
		"StateText": state.StateText[j.State],
	})
}