package workflow

type (
	// Context represents the context of the current workflow request.
	Context interface {
		// Get retrieves data from the context.
		Get(key string) interface{}

		// Set saves data in the context.
		Set(key string, val interface{})

		// FlowParam returns flow parameter by name.
		FlowParam(name string) interface{}

		// SetFlowParam set flow parameter by name.
		SetFlowParam(name string, value interface{}) error

		// TaskParam returns task parameter by name.
		TaskParam(name string) interface{}

		// SetTaskParam set task parameter by name.
		SetTaskParam(name string, value interface{}) error
	}
)
