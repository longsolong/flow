package atom

import (
	"context"
	"fmt"
	"github.com/longsolong/flow/pkg/orchestration/request"
)

// Atom ...
type Atom interface {
	Create(ctx context.Context, req *request.Request) error
	AtomID() AtomID
}

// AtomID ...
type AtomID struct {
	Type            string
	ID              string
	ExpansionDigest string
}


func (id AtomID) String() string {
	return fmt.Sprintf("AtomID: %s.%s.%s", id.Type, id.ID, id.ExpansionDigest)
}