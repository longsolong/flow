package workflow

import (
	"github.com/longsolong/flow/pkg/orchestration/request"
)

type (
	// Context represents the context of the current workflow request.
	Context interface {
		// Request returns flow request
		Request() interface{}

		// SetRequest sets `*http.Request`.
		SetRequest(r *request.Request)

		// FlowParam returns flow parameter by name.
		FlowParam(name string) interface{}

		// SetFlowParam set flow parameter by name.
		SetFlowParam(name string, value interface{}) error

		// JobParam returns job parameter by name.
		JobParam(name string) interface{}

		// SetJobParam set job parameter by name.
		SetJobParam(name string, value interface{}) error
	}
)
