// Code generated by "genatom -type=Generate"; DO NOT EDIT.

package primesieve

import (
	"reflect"

	"github.com/longsolong/flow/pkg/workflow/atom"
)

// AtomID ...
func (s *Generate) AtomID() atom.AtomID {
	return atom.AtomID{
		Type:            reflect.TypeOf(s).Elem().String(),
		ID:              s.ID,
		ExpansionDigest: s.ExpansionDigest,
	}
}
