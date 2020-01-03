package request

import (
	"context"
	"github.com/google/uuid"
)

// Request ...
type Request struct {
	RequestUUID uuid.UUID
	RequestArgs map[string]interface{}
	RequestTags []Tag

	ctx context.Context
}

// NewRequest ...
func NewRequest() *Request {
	return NewRequestWithContext(context.Background())
}

// NewRequestWithContext ...
func NewRequestWithContext(ctx context.Context) *Request {
	req := new(Request)
	req.ctx = ctx
	req.RequestUUID = uuid.New()
	req.RequestTags = make([]Tag, 0)
	return req
}

// Tag ...
type Tag struct {
	Name  string
	Value string
}

// Context ...
func (r *Request) Context() context.Context {
	if r.ctx != nil {
		return r.ctx
	}
	return context.Background()
}

// WithContext returns a shallow copy of r with its context changed
// to ctx. The provided ctx must be non-nil.
func (r *Request) WithContext(ctx context.Context) *Request {
	if ctx == nil {
		panic("nil context")
	}
	r2 := new(Request)
	*r2 = *r
	r2.ctx = ctx
	return r2
}
