package atom

import "github.com/longsolong/flow/pkg/workflow"

// Atom ...
type Atom interface {
	Create(ctx workflow.Context) error
	ID() ID
}

// ID ...
type ID struct {
	Type            string
	ID              string
	ExpansionDigest string
}

