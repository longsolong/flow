package context

import (
	"context"
	"github.com/longsolong/flow/pkg/infra"
)

type FlowContextKey string

var (
	LoggerCtxKey     = FlowContextKey("Logger")
)

// Logger returns the logger from a context object.
func Logger(ctx context.Context) *infra.Logger {
	logger, ok := ctx.Value(LoggerCtxKey).(*infra.Logger)
	if !ok {
		panic("valve: LoggerCtxKey has not been set on the context.")
	}
	return logger
}

