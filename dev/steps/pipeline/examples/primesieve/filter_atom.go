// Code generated by "genatom -type=Filter"; DO NOT EDIT.

package primesieve

import (
	"reflect"

	"github.com/longsolong/flow/pkg/workflow/atom"
)

// AtomID ...
func (s *Filter) AtomID() atom.AtomID {
	return atom.AtomID{
		Type:            reflect.TypeOf(s).Elem().String(),
		ID:              s.ID,
		ExpansionDigest: s.ExpansionDigest,
	}
}
