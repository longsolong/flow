package middleware

import (
	"github.com/longsolong/flow/pkg/infra"
)

// Middleware wrapps all middlewares
type Middleware struct {
	logger *infra.Logger
}

// CreateMiddleware create a instance of a middleware
func CreateMiddleware(l *infra.Logger) *Middleware {
	return &Middleware{
		logger: l,
	}
}
