package rest

import (
	"io"
	"net/http"
)

// NewHealthCheckHandler add route for healthcheck
func (h *Handler) NewHealthCheckHandler() {
	h.router.Get("/health", healthCheckHandler)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}
