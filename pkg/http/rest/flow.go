package rest

import (
	"context"
	"encoding/json"
	"github.com/go-chi/valve"
	"github.com/longsolong/flow/dev/workflows/pipeline"
	"github.com/longsolong/flow/dev/workflows/standalone"
	"io/ioutil"
	"net/http"

	"github.com/faceair/jio"
	"github.com/go-chi/chi"
	"github.com/longsolong/flow/pkg/infra"
	flowcontext "github.com/longsolong/flow/pkg/workflow/context"
)

// NewFlowHandler add route for flow
func (h *Handler) NewFlowHandler(logger *infra.Logger) {
	singleProcessorFlowHandler := SingleProcessorFlowHandler{logger: logger}
	h.router.Route("/api/standalone/flows", func(r chi.Router) {
		r.With(jio.ValidateBody(RunFlowValidator, jio.DefaultErrorHandler)).Post("/run", singleProcessorFlowHandler.Run())
	})
	pipelineFlowHandler := PipelineFlowHandler{logger: logger}
	h.router.Route("/api/pipeline/flows", func(r chi.Router) {
		r.With(jio.ValidateBody(RunFlowValidator, jio.DefaultErrorHandler)).Post("/run", pipelineFlowHandler.Run())
	})
}

// RunFlowValidator ...
var RunFlowValidator = jio.Object().Keys(jio.K{
	"primaryRequestArgs": jio.Object().Keys(jio.K{
		"namespace": jio.String().Required(),
		"name":      jio.String().Required(),
		"version":   jio.Number().Integer().Required(),
	}),
	"requestArgs": jio.Object().Required(),
	"requestTags": jio.Array().Required(),
})

// SingleProcessorFlowHandler ...
type SingleProcessorFlowHandler struct {
	logger   *infra.Logger
}

// Run ...
func (h SingleProcessorFlowHandler) Run() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			jio.DefaultErrorHandler(w, r, err)
			return
		}
		data, err := jio.ValidateJSON(&body, RunFlowValidator)
		if err != nil {
			jio.DefaultErrorHandler(w, r, err)
			return
		}
		valv := valve.New()
		baseCtx := valv.Context()
		r = r.WithContext(baseCtx)
		namespace := data["primaryRequestArgs"].(map[string]interface{})["namespace"]
		name := data["primaryRequestArgs"].(map[string]interface{})["name"]
		version := data["primaryRequestArgs"].(map[string]interface{})["version"]
		grapher, err := standalone.SingleProcessorFactory.Make(
			context.WithValue(r.Context(), flowcontext.LoggerCtxKey, h.logger),
			h.logger, namespace.(string), name.(string), int(version.(float64)), body)
		if err != nil {
			jio.DefaultErrorHandler(w, r, err)
			return
		}
		b, err := json.Marshal(grapher.Chain.AllJobs())
		if err != nil {
			jio.DefaultErrorHandler(w, r, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		return
	}
	return fn
}

// PipelineFlowHandler ...
type PipelineFlowHandler struct {
	logger   *infra.Logger
}

// Run ...
func (h PipelineFlowHandler) Run() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			jio.DefaultErrorHandler(w, r, err)
			return
		}
		data, err := jio.ValidateJSON(&body, RunFlowValidator)
		if err != nil {
			jio.DefaultErrorHandler(w, r, err)
			return
		}
		valv := valve.New()
		baseCtx := valv.Context()
		r = r.WithContext(baseCtx)
		namespace := data["primaryRequestArgs"].(map[string]interface{})["namespace"]
		name := data["primaryRequestArgs"].(map[string]interface{})["name"]
		version := data["primaryRequestArgs"].(map[string]interface{})["version"]
		_, err = pipeline.PipelineFactory.Make(
			context.WithValue(r.Context(), flowcontext.LoggerCtxKey, h.logger),
			h.logger, namespace.(string), name.(string), int(version.(float64)), body)
		if err != nil {
			jio.DefaultErrorHandler(w, r, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		return
	}
	return fn
}