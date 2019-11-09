package step

import (
	"github.com/longsolong/flow/pkg/workflow"
)

// Atom ...
type Atom interface {
	Create(ctx workflow.Context) error
}

// ID ...
type ID struct {
	Type            string
	Name            string
	ID              string
	ExpansionDigest string
}

// Step ...
type Step struct {
	ID         ID
	SequenceID *ID
	Atom
}
