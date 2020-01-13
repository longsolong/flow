package workflow

import (
	"errors"
)

var (
	// ErrAlreadyRegisteredNode ...
	ErrAlreadyRegisteredNode = errors.New("already registered node")

	// ErrNotRegisteredNode ...
	ErrNotRegisteredNode = errors.New("not registered node")

	// ErrAlreadyRegisteredUpstream ...
	ErrAlreadyRegisteredUpstream = errors.New("already registered upstream")

	// ErrAlreadyRegisteredDownstream ...
	ErrAlreadyRegisteredDownstream = errors.New("already registered downstream")
)
