// Code generated by "genatom -type=Noop"; DO NOT EDIT.

package builtin

import (
	"reflect"

	"github.com/longsolong/flow/pkg/workflow/atom"
)

// AtomID ...
func (s *Noop) AtomID() atom.AtomID {
	return atom.AtomID{
		Type:            reflect.TypeOf(s).Elem().String(),
		ID:              s.ID,
		ExpansionDigest: s.ExpansionDigest,
	}
}
