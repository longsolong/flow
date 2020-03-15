package pipeline

import "context"

// Pipeline ...
type Pipeline interface {
	Run(ctx context.Context) error
	Stop(ctx context.Context) error
}
