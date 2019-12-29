package infra

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/valve"
)


// CreateServer create a new http server
func CreateServer(router chi.Router, serverAddress string, options ...func(*http.Server)) *http.Server {
	// Our graceful valve shut-off package to manage code preemption and
	// shutdown signaling.
	valv := valve.New()
	baseCtx := valv.Context()

	server := &http.Server{
		Addr:         serverAddress,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      chi.ServerBaseContext(baseCtx, router),
	}
	for _, option := range options {
		option(server)
	}
	return server
}
