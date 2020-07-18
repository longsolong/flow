package numberguess

import (
	"context"
	"time"

	"github.com/faceair/jio"
	"github.com/longsolong/flow/dev/steps/standalone/examples/numberguess"
	"github.com/longsolong/flow/pkg/orchestration/request"
	"github.com/longsolong/flow/pkg/orchestration/standalone/graph"
)

const (
	// NAME ...
	NAME = "number_guess"
	// VERSION ...
	VERSION = 1
)

var schema = jio.Object().Keys(jio.K{
	"requestArgs": jio.Object().Keys(jio.K{
		"secret": jio.Number().Integer().Min(numberguess.LOW).Max(numberguess.HIGH).Required(),
	}),
	"requestTags": jio.Array().Items(jio.Object().Keys(jio.K{
		"name":  jio.String().Required(),
		"value": jio.String().Required(),
	})),
})

//go:generate gengrapher -type=NumberGuess

type plotter struct {
	graph.Plotter
}


//go:generate genaccessor -type=numberGuessParam

type numberGuessParam struct {
	Secret int
	Low    int
	High   int
}

func (p *plotter) Begin(ctx context.Context, req *request.Request) error {
	if _, err := p.NewNode(
		ctx, req,
		numberguess.NewNumberGuess("", "", &numberGuessParam{}),
		"try guess number by binary search", 6, time.Duration(10)*time.Millisecond); err != nil {
		return err
	}
	return nil
}

func (p *plotter) Grow(ctx context.Context) {
	p.Plotter.Close()
}
