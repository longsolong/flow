package rest

import (
	"github.com/faceair/jio"
	"github.com/go-chi/chi"
	"github.com/longsolong/flow/pkg/http/rest/single_processor/flow"
)

// NewFlowHandler add route for flow
func (h *Handler) NewFlowHandler() {
	h.router.Route("/api/single_processor/flows", func(r chi.Router) {
		r.With(jio.ValidateBody(flow.RunFlowValidator, jio.DefaultErrorHandler)).Post("/run", flow.RunFlowHandler)
	})
}
