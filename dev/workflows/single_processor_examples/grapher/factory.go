package grapher

import (
	"fmt"
	"github.com/longsolong/flow/pkg/orchestration/single_processor/grapher"
	"github.com/longsolong/flow/dev/workflows/single_processor_examples/ping"
)


// grapherFactory implements the Factory interface.
type grapherFactory struct {
}

// NewGrapherFactory creates a Factory.
func NewGrapherFactory() grapher.Factory {
	return &grapherFactory{
	}
}

// Make
func (gf *grapherFactory) Make(name string, version int, rawRequestArgs []byte) (*grapher.Grapher, error) {
	if name == "ping" && version == 1 {
		return ping.NewGrapher(name, version, rawRequestArgs)
	}
	return nil, fmt.Errorf("unknown dag name %s version %d", name, version)
}

