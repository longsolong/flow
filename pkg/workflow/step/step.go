package step

import (
	"github.com/longsolong/flow/pkg/workflow"
)

type Atom interface {
	Create(ctx workflow.Context) error
}

type Step struct {
	Type string
	Description string
	Atom
}