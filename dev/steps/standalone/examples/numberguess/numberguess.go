package numberguess

import (
	"context"
	"fmt"
	"github.com/longsolong/flow/pkg/orchestration/request"
	"github.com/longsolong/flow/pkg/workflow/atom"
	flowcontext "github.com/longsolong/flow/pkg/workflow/context"
	"github.com/longsolong/flow/pkg/workflow/state"
	"github.com/longsolong/flow/pkg/workflow/step"
)

const (
	LOW = 0
	HIGH = 100
)

//go:generate genatom -type=NumberGuess

// NumberGuess ...
type NumberGuess struct {
	*step.Step
	NumberGuessParameter
}

// NumberGuessParameter ...
type NumberGuessParameter interface {
	MustSetSecret(secret int)
	GetSecret() int
	MustSetLow(low int)
	GetLow() int
	MustSetHigh(high int)
	GetHigh() int
}

// NewNumberGuess ...
func NewNumberGuess(id, expansionDigest string, parameter NumberGuessParameter) *NumberGuess {
	return &NumberGuess{
		step.NewStep(id, expansionDigest),
		parameter,
	}
}

// Create ...
func (s *NumberGuess) Create(ctx context.Context, req *request.Request) error {
	logger := flowcontext.Logger(ctx)
	secret := int(req.RequestArgs["secret"].(float64))
	s.MustSetSecret(secret)
	s.MustSetLow(LOW)
	s.MustSetHigh(HIGH)
	logger.Log.Info(fmt.Sprintf("secret is %d", secret))
	return nil
}

// Run ...
func (s *NumberGuess) Run(ctx context.Context) (atom.Return, error) {
	ret := atom.Return{}
	secret := s.GetSecret()
	low := s.GetLow()
	high := s.GetHigh()
	logger := flowcontext.Logger(ctx)
	totalTries := ctx.Value(flowcontext.FlowContextKey("totalTries")).(uint)
	guess := (low + high) / 2
	logger.Log.Info(fmt.Sprintf("%d guess is %d", totalTries, guess))
	if guess != secret {
		ret.Exit = 1
		ret.State = state.StateFail
		if guess < secret {
			s.MustSetLow(guess+1)
		} else {
			s.MustSetHigh(guess-1)
		}
	}
	return ret, nil
}

// Stop run
func (s *NumberGuess) Stop(ctx context.Context) error {
	return nil
}
