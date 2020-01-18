package workflows

import (
	"context"
	"fmt"
	"time"

	"github.com/longsolong/flow/dev/workflows/single_processor_examples/ping"
	"github.com/longsolong/flow/pkg/execution/single_processor/traverser"
	"github.com/longsolong/flow/pkg/infra"
	"github.com/longsolong/flow/pkg/orchestration/single_processor/graph"
)

// SingleProcessorFactory ...
var SingleProcessorFactory = singleProcessorFactory{}

// singleProcessorFactory ...
type singleProcessorFactory struct{}

// Make
func (gf *singleProcessorFactory) Make(ctx context.Context, logger *infra.Logger, namespace, name string, version int, rawRequestData []byte) (g *graph.Grapher, err error) {
	if namespace == "single_processor_examples" {
		switch {
		case name == ping.NAME && version == ping.VERSION:
			g, err = ping.NewGrapher(ctx, rawRequestData)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unknown dag name %s version %d", name, version)
		}
	} else {
		return nil, fmt.Errorf("unknown namespace %s", namespace)
	}
	t := traverser.NewTraverser(g, logger, time.Duration(10)*time.Second, time.Duration(10)*time.Second)
	go g.GraphPlotter.Grow(ctx)
	t.Run(ctx)

	return g, nil
}
