package standalone

import (
	"context"
	"fmt"
	"github.com/longsolong/flow/dev/workflows/standalone/examples/numberguess"
	"time"

	"github.com/longsolong/flow/pkg/execution/standalone/traverser"
	"github.com/longsolong/flow/pkg/infra"
	"github.com/longsolong/flow/pkg/orchestration/standalone/graph"
)

// SingleProcessorFactory ...
var SingleProcessorFactory = singleProcessorFactory{}

// singleProcessorFactory ...
type singleProcessorFactory struct{}

// Make
func (gf *singleProcessorFactory) Make(ctx context.Context, logger *infra.Logger, namespace, name string, version int, rawRequestData []byte) (g *graph.Grapher, err error) {
	if namespace == "examples" {
		switch {
		case name == numberguess.NAME && version == numberguess.VERSION:
			g, err = numberguess.NewGrapher(ctx, rawRequestData)
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
