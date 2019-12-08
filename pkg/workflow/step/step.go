package step

import (
	"github.com/longsolong/flow/pkg/workflow/atom"
)

// Step ...
type Step struct {
	id atom.ID
}

// ID ...
func(s *Step) ID() atom.ID {
	return s.id
}

// SetID ...
func(s *Step) SetID(id atom.ID) {
	s.id = id
}