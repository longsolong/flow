package primesieve

import (
	"context"
	"fmt"
	"github.com/go-chi/valve"
	"github.com/longsolong/flow/pkg/orchestration/request"
	"github.com/longsolong/flow/pkg/workflow/atom"
	flowcontext "github.com/longsolong/flow/pkg/workflow/context"
	"github.com/longsolong/flow/pkg/workflow/step"
)

//go:generate genatom -type=Generate

// Generate ...
type Generate struct {
	*step.Step
	GenerateParameter
}

// GenerateParameter ...
type GenerateParameter interface {
	MustSetSeqCh(chan<- int)
	GetSeqCh() chan<- int
}

// NewGenerate ...
func NewGenerate(id, expansionDigest string, parameter GenerateParameter) *Generate {
	return &Generate{
		Step: step.NewStep(id, expansionDigest),
		GenerateParameter: parameter,
	}
}

// Create ...
func (s *Generate) Create(ctx context.Context, req *request.Request) error {
	return nil
}

// Run ...
func (s *Generate) Run(ctx context.Context) (atom.Return, error) {
	ret := atom.Return{}
	logger := flowcontext.Logger(ctx)
	valv := valve.Lever(ctx)
	for i := 2; ; i++ {
		select {
		case s.GetSeqCh() <- i:
			logger.Log.Info(fmt.Sprintf("Send '%d' to channel 'SeqCh'.", i))
		case <-valv.Stop():
			break
		}
	}
	return ret, nil
}

// Stop run
func (s *Generate) Stop(ctx context.Context) error {
	return nil
}

//go:generate genatom -type=Filter

// Filter ...
type Filter struct {
	*step.Step
	FilterParameter
}

type FilterParameter interface {
	MustSetInCh(<-chan int)
	GetInCh() <-chan int
	MustSetOutCh(chan<- int)
	GetOutCh() chan<- int
	MustSetPrime(prime int)
	GetPrime() int
}

// NewFilter ...
func NewFilter(id, expansionDigest string, parameter FilterParameter) *Filter {
	return &Filter{
		step.NewStep(id, expansionDigest),
		parameter,
	}
}

// Create ...
func (s *Filter) Create(ctx context.Context, req *request.Request) error {
	return nil
}

// Run ...
func (s *Filter) Run(ctx context.Context) (atom.Return, error) {
	ret := atom.Return{}
	logger := flowcontext.Logger(ctx)
	valv := valve.Lever(ctx)
	for {
		select {
		case i := <-s.GetInCh():
			logger.Log.Info(fmt.Sprintf("Receive value %d from InCh.", i))
			if i%s.GetPrime() != 0 {
				s.GetOutCh() <- i
				logger.Log.Info(fmt.Sprintf("Send %d to OutCh.", i))
			}
		case <-valv.Stop():
			break
		}
	}
	return ret, nil
}

// Stop run
func (s *Filter) Stop(ctx context.Context) error {
	return nil
}
