// Code generated by "genatom -type=NumberGuess"; DO NOT EDIT.

package numberguess

import (
	"reflect"

	"github.com/longsolong/flow/pkg/workflow/atom"
)

// AtomID ...
func (s *NumberGuess) AtomID() atom.AtomID {
	return atom.AtomID{
		Type:            reflect.TypeOf(s).Elem().String(),
		ID:              s.ID,
		ExpansionDigest: s.ExpansionDigest,
	}
}