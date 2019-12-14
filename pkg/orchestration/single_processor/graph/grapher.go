package graph

import (
	"github.com/longsolong/flow/pkg/orchestration/request"

	"github.com/longsolong/flow/pkg/orchestration/single_processor/chain"
	"github.com/longsolong/flow/pkg/workflow/dag"
)

type (
	buildRequestFunc            func(rawRequestData []byte) (*request.Request, error)
)

// Grapher ...
type Grapher struct {
	req   *request.Request
	dag   *dag.DAG
	chain *chain.Chain
	plotter Plotter

	buildRequest buildRequestFunc
}

// NewGrapher A new Grapher should be made for every request.
func NewGrapher(
	name string, version int,
	rawRequestData []byte,
	buildRequest buildRequestFunc,
	plotter Plotter) (*Grapher, error) {

	g := &Grapher{
		buildRequest: buildRequest,
		plotter: plotter,
	}
	req, err := (g.buildRequest)(rawRequestData)
	if err != nil {
		return nil, err
	}
	g.req = req

	d, c, err := g.plotter.Begin(name, version, req)
	if err != nil {
		return nil, err
	}
	g.dag = d
	g.chain = c

	return g, nil
}
