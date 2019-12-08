package workflow

import (
	"errors"
)

var (
	// ErrAlreadyRegisteredNode ...
	ErrAlreadyRegisteredNode = errors.New("already registered node")

	// ErrAlreadyRegisteredUpstream ...
	ErrAlreadyRegisteredUpstream = errors.New("already registered upstream")

	// ErrAlreadyRegisteredDownstream ...
	ErrAlreadyRegisteredDownstream = errors.New("already registered downstream")
)
