package atom

import (
	"context"
	"github.com/longsolong/flow/pkg/orchestration/request"
)

// Atom ...
type Atom interface {
	Create(ctx context.Context, req *request.Request) error
	ID() ID
}

// ID ...
type ID struct {
	Type            string
	ID              string
	ExpansionDigest string
}

