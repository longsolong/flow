package job

import (
	"github.com/google/uuid"
	"github.com/longsolong/flow/pkg/workflow/atom"
	"github.com/longsolong/flow/pkg/workflow/state"
)

// Job ...
type Job struct {
	atom.Atom             // step
	State     state.State // State const

	// RequestId of the request that created the job. This is only informational
	// for reporting/loggging/tracing.
	RequestUUID uuid.UUID
}

// NewJob ...
func NewJob(step atom.Atom, requestUUID uuid.UUID) *Job {
	return &Job{
		Atom:        step,
		State:       state.StateUnknown,
		RequestUUID: requestUUID,
	}
}
