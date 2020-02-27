package workflows

import (
	"context"
	"github.com/go-chi/valve"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/longsolong/flow/pkg/infra"
	flowcontext "github.com/longsolong/flow/pkg/workflow/context"
)

func TestNewDAG(t *testing.T) {
	logger, err := infra.CreateLogger(0)
	if err != nil {
		t.Fatal(err)
	}

	namespace := "single_processor_examples"
	name := "number_guess"
	version := 1

	valv := valve.New()
	baseCtx := valv.Context()
	body := []byte(`{
		"primaryRequestArgs": {
			"namespace": "single_processor_examples",
			"name": "number_guess",
			"version": 1
		},
		"requestArgs": {
			"secret": 1
		},
		"requestTags": [
			{"name": "aa", "value": "bb"}
		]
	}`)
	_, err = SingleProcessorFactory.Make(
		context.WithValue(baseCtx, flowcontext.FlowContextKey("logger"), logger),
		logger, namespace, name, version, body)
	assert.Nil(t, err)
}
