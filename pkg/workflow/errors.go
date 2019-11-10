package workflow

import (
	"errors"
)

var (
	// ErrAlreadyRegisteredStep ...
	ErrAlreadyRegisteredStep = errors.New("already registered step")

	// ErrAlreadyRegisteredUpstream ...
	ErrAlreadyRegisteredUpstream = errors.New("already registered upstream")
)
