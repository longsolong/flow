package pipeline

import (
	"context"
	"github.com/go-chi/valve"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/longsolong/flow/pkg/infra"
	flowcontext "github.com/longsolong/flow/pkg/workflow/context"
)

func TestNewPipeline(t *testing.T) {
	logger, err := infra.CreateLogger(0)
	if err != nil {
		t.Fatal(err)
	}

	namespace := "examples"
	name := "prime_sieve"
	version := 1

	valv := valve.New()
	baseCtx := valv.Context()
	body := []byte(`{
		"primaryRequestArgs": {
			"namespace": "examples",
			"name": "prime_sieve",
			"version": 1
		},
		"requestArgs": {
			"num": 10
		},
		"requestTags": [
			{"name": "aa", "value": "bb"}
		]
	}`)
	_, err = PipelineFactory.Make(
		context.WithValue(baseCtx, flowcontext.LoggerCtxKey, logger),
		logger, namespace, name, version, body)
	valv.Shutdown(1)
	assert.Nil(t, err)
}
