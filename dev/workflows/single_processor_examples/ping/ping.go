package ping

import (
	"github.com/faceair/jio"
	"github.com/longsolong/flow/dev/steps/examples"
	"github.com/longsolong/flow/pkg/orchestration/request"
	"github.com/longsolong/flow/pkg/orchestration/single_processor/chain"
	"github.com/longsolong/flow/pkg/orchestration/single_processor/grapher"
	"github.com/longsolong/flow/pkg/workflow/dag"
	"time"
)

func buildRequest(rawRequestArgs []byte) (*request.Request, error) {
	requestArgsSchema := jio.Object().Keys(jio.K{
		"args": jio.Object().Keys(jio.K{
			"hostname": jio.String().Required(),
			"timeout":  jio.Number().Integer().Required(),
			"interval": jio.Number().Integer().Required(),
			"count":    jio.Number().Integer().Required(),
		}),
		"tags": jio.Array().Items(jio.Object().Keys(jio.K{
			"name":  jio.String().Required(),
			"value": jio.String().Required(),
		})),
	})
	requestArgs, err := jio.ValidateJSON(&rawRequestArgs, requestArgsSchema)
	if err != nil {
		return nil, err
	}
	req := request.NewRequest()
	req.PrimaryRequestArgs = requestArgs["args"].(map[string]interface{})
	for _, v := range requestArgs["tags"].([]map[string]interface{}) {
		for kk, vv := range v {
			req.PrimaryRequestTags = append(req.PrimaryRequestTags, request.Tag{Name: kk, Value: vv.(string)})
		}
	}
	return req, nil
}

func buildStepDag(name string, version int, req *request.Request) (*dag.DAG, error) {
	d := dag.NewDAG(name, version)
	step1 := dag.NewNode(examples.NewPing("", ""), "ping host", 3, time.Duration(10)*time.Millisecond)
	err := d.AddNode(step1)

	return d, err
}

func buildJobChain(req *request.Request, d *dag.DAG) (*chain.Chain, error) {
	c := chain.NewChain(d, req.RequestUUID)
	return c, nil
}

// NewGrapher ...
func NewGrapher(name string, version int, rawRequestArgs []byte) (*grapher.Grapher, error) {
	return grapher.NewGrapher(
		name, version,
		rawRequestArgs,
		buildRequest, buildStepDag, buildJobChain, nil, nil)
}
