package grapher

import (
	"github.com/longsolong/flow/pkg/orchestration/request"

	"github.com/longsolong/flow/pkg/orchestration/single_processor/chain"
	"github.com/longsolong/flow/pkg/workflow/dag"
)

type (
	buildRequestFunc            func(rawRequestArgs []byte) (*request.Request, error)
	buildStepDagFunc            func(name string, version int, req *request.Request) (*dag.DAG, error)
	buildJobChainFunc           func(*request.Request, *dag.DAG) (*chain.Chain, error)
	dynamicPopulateStepDagFunc  func(*dag.DAG, *request.Request) error
	dynamicPopulateJobChainFunc func(*chain.Chain, *request.Request) error
)

// Grapher ...
type Grapher struct {
	req   *request.Request
	dag   *dag.DAG
	chain *chain.Chain

	BuildRequest            buildRequestFunc
	BuildStepDag            buildStepDagFunc
	BuildJobChain           buildJobChainFunc
	DynamicPopulateStepDag  *dynamicPopulateStepDagFunc
	DynamicPopulateJobChain *dynamicPopulateJobChainFunc
}

// NewGrapher A new Grapher should be made for every request.
func NewGrapher(
	name string, version int,
	rawRequestArgs []byte,
	buildRequest buildRequestFunc,
	buildStepDag buildStepDagFunc,
	buildJobChain buildJobChainFunc,
	dynamicPopulateStepDag *dynamicPopulateStepDagFunc,
	dynamicPopulateJobChain *dynamicPopulateJobChainFunc) (*Grapher, error) {

	g := &Grapher{
		BuildRequest:            buildRequest,
		BuildStepDag:            buildStepDag,
		BuildJobChain:           buildJobChain,
		DynamicPopulateStepDag:  dynamicPopulateStepDag,
		DynamicPopulateJobChain: dynamicPopulateJobChain,
	}
	req, err := (g.BuildRequest)(rawRequestArgs)
	if err != nil {
		return nil, err
	}
	g.req = req

	d, err := (g.BuildStepDag)(name, version, req)
	if err != nil {
		return nil, err
	}
	g.dag = d

	c, err := (g.BuildJobChain)(g.req, d)
	if err != nil {
		return nil, err
	}
	g.chain = c

	return g, nil
}
