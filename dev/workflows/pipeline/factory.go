package pipeline

import (
	"context"
	"fmt"
	"github.com/longsolong/flow/dev/workflows/pipeline/examples/primesieve"
	pipe "github.com/longsolong/flow/pkg/execution/pipeline"

	"github.com/longsolong/flow/pkg/infra"
)

// PipelineFactory ...
var PipelineFactory = pipelineFactory{}

// pipelineFactory ...
type pipelineFactory struct{}

// Make
func (gf *pipelineFactory) Make(ctx context.Context, logger *infra.Logger, namespace, name string, version int, rawRequestData []byte) (p pipe.Pipeline, err error) {
	if namespace == "examples" {
		switch {
		case name == primesieve.NAME && version == primesieve.VERSION:
			p, err = primesieve.NewPipeline(ctx, rawRequestData)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unknown dag name %s version %d", name, version)
		}
	} else {
		return nil, fmt.Errorf("unknown namespace %s", namespace)
	}
	err = p.Run(ctx)

	return p, err
}
