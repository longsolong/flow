package ping

import (
	"context"
	"time"

	"github.com/faceair/jio"
	"github.com/longsolong/flow/dev/steps/examples"
	"github.com/longsolong/flow/pkg/orchestration/job"
	"github.com/longsolong/flow/pkg/orchestration/request"
	"github.com/longsolong/flow/pkg/orchestration/single_processor/graph"
	"github.com/longsolong/flow/pkg/workflow/dag"
)

const (
	// NAME ...
	NAME = "ping"
	// VERSION ...
	VERSION = 1
)

type plotter struct {
	graph.Plotter
}

func (p *plotter) Begin(ctx context.Context, req *request.Request) error {
	// node
	step1 := examples.NewPing("", "")
	if err := step1.Create(ctx, req); err != nil {
		return err
	}
	err := p.DAG.AddNode(dag.NewNode(step1, "ping host", 3, time.Duration(10)*time.Millisecond))
	if err != nil {
		return err
	}

	// job
	p.Chain.AddJob(job.NewJob(step1))

	return nil
}

func (p *plotter) Grow(ctx context.Context) {
	p.Plotter.Close()
}

// NewGrapher ...
func NewGrapher(ctx context.Context, rawRequestData []byte) (*graph.Grapher, error) {
	req, err := newRequest(ctx, rawRequestData)
	if err != nil {
		return nil, err
	}
	p, err := newPlotter(ctx, req)
	if err != nil {
		return nil, err
	}
	g := graph.NewGrapher(req, p.DAG, p.Chain, p)
	return g, nil
}

func newRequest(ctx context.Context, rawRequestData []byte) (*request.Request, error) {
	schema := jio.Object().Keys(jio.K{
		"requestArgs": jio.Object().Keys(jio.K{
			"hostname": jio.String().Required(),
			"timeout":  jio.Number().Integer().Required(),
			"interval": jio.Number().Integer().Required(),
			"count":    jio.Number().Integer().Required(),
		}),
		"requestTags": jio.Array().Items(jio.Object().Keys(jio.K{
			"name":  jio.String().Required(),
			"value": jio.String().Required(),
		})),
	})
	requestArgs, err := jio.ValidateJSON(&rawRequestData, schema)
	if err != nil {
		return nil, err
	}
	req := request.NewRequestWithContext(ctx)
	req.RequestArgs = requestArgs["requestArgs"].(map[string]interface{})
	for _, v := range requestArgs["requestTags"].([]interface{}) {
		v := v.(map[string]interface{})
		req.RequestTags = append(req.RequestTags, request.Tag{Name: v["name"].(string), Value: v["value"].(string)})
	}
	return req, nil
}

// newPlotter ...
func newPlotter(ctx context.Context, req *request.Request) (*plotter, error) {
	p := &plotter{Plotter: graph.NewPlotter(NAME, VERSION)}
	if err := p.Begin(ctx, req); err != nil {
		return nil, err
	}
	return p, nil
}
