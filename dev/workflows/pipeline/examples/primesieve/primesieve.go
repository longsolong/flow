package primesieve

import (
	"context"
	"fmt"
	"github.com/go-chi/valve"
	"github.com/longsolong/flow/dev/steps/pipeline/examples/primesieve"
	flowcontext "github.com/longsolong/flow/pkg/workflow/context"

	"github.com/faceair/jio"
	pipe "github.com/longsolong/flow/pkg/execution/pipeline"
	"github.com/longsolong/flow/pkg/orchestration/request"
)

// https://play.golang.org/p/9U22NfrXeq

const (
	// NAME ...
	NAME = "prime_sieve"
	// VERSION ...
	VERSION = 1
)

//go:generate genaccessor -type=generateParam,filterParam -v

type generateParam struct {
	SeqCh chan<- int
}

type filterParam struct {
	InCh  <-chan int
	OutCh chan<- int
	Prime int
}

// pipeline ...
type pipeline struct {
	Req *request.Request
}

var schema = jio.Object().Keys(jio.K{
	"requestArgs": jio.Object().Keys(jio.K{
		"num": jio.Number().Integer().Required(),
	}),
	"requestTags": jio.Array().Items(jio.Object().Keys(jio.K{
		"name":  jio.String().Required(),
		"value": jio.String().Required(),
	})),
})

// NewPipeline ...
func NewPipeline(ctx context.Context, rawRequestData []byte) (pipe.Pipeline, error) {
	req, err := newRequest(ctx, rawRequestData)
	if err != nil {
		return nil, err
	}
	return &pipeline{Req: req}, nil
}

func newRequest(ctx context.Context, rawRequestData []byte) (*request.Request, error) {
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

// Run ...
func (p *pipeline) Run(ctx context.Context) error {
	logger := flowcontext.Logger(ctx)
	ch := make(chan int) // Create a new channel.
	g := primesieve.NewGenerate("", "", &generateParam{SeqCh: ch})
	go g.Run(ctx) // Launch Generate goroutine.
	for i := 0; i < 10; i++ {
		prime := <-ch
		logger.Log.Info(fmt.Sprintf("Prime %d.", prime))
		ch1 := make(chan int)
		f := primesieve.NewFilter("", "", &filterParam{InCh: ch, OutCh: ch1, Prime: prime})
		go f.Run(ctx)
		ch = ch1
	}
	valv := valve.Lever(ctx)
	valv.(*valve.Valve).Shutdown(0)
	return nil
}

// Stop ...
func (p *pipeline) Stop(ctx context.Context) error {
	return nil
}
