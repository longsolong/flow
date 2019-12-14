package ping

import (
	"github.com/faceair/jio"
	"github.com/longsolong/flow/dev/steps/examples"
	"github.com/longsolong/flow/pkg/orchestration/job"
	"github.com/longsolong/flow/pkg/orchestration/request"
	"github.com/longsolong/flow/pkg/orchestration/single_processor/chain"
	"github.com/longsolong/flow/pkg/orchestration/single_processor/graph"
	"github.com/longsolong/flow/pkg/workflow/dag"
	"time"
)

const (
	// NAME ...
	NAME = "ping"
	// VERSION ...
	VERSION = 1
)

func buildRequest(rawRequestData []byte) (*request.Request, error) {
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
	req := request.NewRequest()
	req.RequestArgs = requestArgs["requestArgs"].(map[string]interface{})
	for _, v := range requestArgs["requestTags"].([]interface{}) {
		v := v.(map[string]interface{})
		for kk, vv := range v {
			vv := vv.(string)
			req.RequestTags = append(req.RequestTags, request.Tag{Name: kk, Value: vv})
		}
	}
	return req, nil
}

type plotter struct {}

func (p *plotter) Begin(name string, version int, req *request.Request) (*dag.DAG, *chain.Chain, error) {
	d := dag.NewDAG(name, version)
	step1 := dag.NewNode(examples.NewPing("", ""), "ping host", 3, time.Duration(10)*time.Millisecond)
	err := d.AddNode(step1)
	if err != nil {
		return nil, nil, err
	}
	c := chain.NewChain(d)
	c.AddJob(job.NewJob(step1.Datum))
	return d, c, err
}

func (p *plotter) Grow(*dag.DAG, *chain.Chain, *request.Request) error {
	return nil
}

// NewGrapher ...
func NewGrapher(rawRequestData []byte) (*graph.Grapher, error) {
	return graph.NewGrapher(
		NAME, VERSION,
		rawRequestData,
		buildRequest,
		new(plotter))
}
