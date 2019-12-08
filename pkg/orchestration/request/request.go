package request

import "github.com/google/uuid"

// Request ...
type Request struct {
	RequestUUID uuid.UUID
	PrimaryRequestArgs map[string]interface{}
	PrimaryRequestTags []Tag
}

// NewRequest ...
func NewRequest() *Request {
	req := new(Request)
	req.RequestUUID = uuid.New()
	req.PrimaryRequestTags = make([]Tag, 0)
	return req
}

// Tag ...
type Tag struct {
	Name string
	Value string
}