package workflows

import (
	"fmt"
	"github.com/longsolong/flow/dev/workflows/single_processor_examples/ping"
	"github.com/longsolong/flow/pkg/orchestration/single_processor/graph"
)

// SingleProcessorFactory ...
var SingleProcessorFactory = singleProcessorFactory{}

// singleProcessorFactory implements the SingleProcessorFactory interface.
type singleProcessorFactory struct {}

// Make
func (gf *singleProcessorFactory) Make(namespace, name string, version int, rawRequestData []byte) (*graph.Grapher, error) {
	if namespace == "single_processor_examples" {
		if name == ping.NAME && version == ping.VERSION {
			return ping.NewGrapher(rawRequestData)
		}
		return nil, fmt.Errorf("unknown dag name %s version %d", name, version)
	}
	return nil, fmt.Errorf("unknown namespace %s", namespace)
}

