package graph

import (
	"context"
	"github.com/longsolong/flow/pkg/orchestration/request"
)

// GraphPlotter ...
type GraphPlotter interface {
	Begin(ctx context.Context, req *request.Request) error
	Grow(ctx context.Context)
	Done() <-chan struct{}
}

