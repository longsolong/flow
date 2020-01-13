package rest

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/faceair/jio"
	"github.com/go-chi/chi"
	"github.com/longsolong/flow/dev/workflows"
	"github.com/longsolong/flow/pkg/infra"
	flowcontext "github.com/longsolong/flow/pkg/workflow/context"
)

// NewFlowHandler add route for flow
func (h *Handler) NewFlowHandler(logger *infra.Logger) {
	singleProcessorFlowHandler := SingleProcessorFlowHandler{logger: logger}
	h.router.Route("/api/single_processor/flows", func(r chi.Router) {
		r.With(jio.ValidateBody(RunFlowValidator, jio.DefaultErrorHandler)).Post("/run", singleProcessorFlowHandler.Run())
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
		namespace := data["primaryRequestArgs"].(map[string]interface{})["namespace"]
		name := data["primaryRequestArgs"].(map[string]interface{})["name"]
		version := data["primaryRequestArgs"].(map[string]interface{})["version"]
		_, err = workflows.SingleProcessorFactory.Make(
			context.WithValue(r.Context(), flowcontext.FlowContextKey("logger"), h.logger),
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
