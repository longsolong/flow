package atom

import (
	"context"
	"fmt"
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


func (id ID) String() string {
	return fmt.Sprintf("Type: %s ID: %s ExpansionDigest: %s", id.Type, id.ID, id.ExpansionDigest)
}