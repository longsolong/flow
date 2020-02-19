package builtin

import (
	"context"
	"github.com/longsolong/flow/pkg/orchestration/request"
	"github.com/longsolong/flow/pkg/workflow/step"
	"reflect"

	"github.com/longsolong/flow/pkg/workflow/atom"
)

// Noop ...
type Noop struct {
	step.Step
}

// NewNoop ...
func NewNoop(id, expansionDigest string) *Noop {
	n := &Noop{}
	n.ID = id
	n.ExpansionDigest = expansionDigest
	return n
}

// Create ...
func (s *Noop) Create(ctx context.Context, req *request.Request) error {
	return nil
}

// Run ...
func (s *Noop) Run(ctx context.Context) (atom.Return, error) {
	ret := atom.Return{}
	return ret, nil
}

// Stop run
func (s *Noop) Stop(ctx context.Context) error {
	return nil
}

// GenAtomType ...
func (s *Noop) GenAtomType() string {
	e := reflect.TypeOf(s).Elem()
	return e.String()
}

// Stop run
func (s *Noop) AtomID() atom.AtomID {
	return atom.AtomID{
		Type: s.GenAtomType(),
		ID: s.ID,
		ExpansionDigest: s.ExpansionDigest,
	}
}
