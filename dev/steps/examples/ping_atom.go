// Code generated by "genatom -type=Ping"; DO NOT EDIT.

package examples

import (
	"reflect"

	"github.com/longsolong/flow/pkg/workflow/atom"
)

// AtomID ...
func (s *Ping) AtomID() atom.AtomID {
	return atom.AtomID{
		Type:            reflect.TypeOf(s).Elem().String(),
		ID:              s.ID,
		ExpansionDigest: s.ExpansionDigest,
	}
}