package rest

import "github.com/longsolong/flow/pkg/http/rest/single_processor/flow"

// NewFlowHandler add route for flow
func (h *Handler) NewFlowHandler() {
	h.router.Post("/single_processor/flows/create_and_run", flow.CreateFlow)
}
