package step

import (
	"github.com/longsolong/flow/pkg/workflow/atom"
)

// Step ...
type Step struct {
	atom.AtomID
}

func (s *Step) StepID() atom.AtomID {
	return s.AtomID
}