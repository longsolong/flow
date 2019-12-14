package rest

import (
	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/longsolong/flow/pkg/http/rest/middleware"
	"github.com/longsolong/flow/pkg/infra"
)

// Handler handles http rest requests
type Handler struct {
	logger *infra.Logger
	router chi.Router
}

// HTTPError data model for http error
type HTTPError struct {
	ErrorCode   int    `json:"error_code"`
	Message     string `json:"message"`
	UserMessage string `json:"user_message"`
}

// CreateHandler create a new http rest handler
func CreateHandler(l *infra.Logger) *Handler {
	h := &Handler{
		logger: l,
		router: chi.NewRouter(),
	}

	// A good base middleware stack
	h.router.Use(chimiddleware.RequestID)
	h.router.Use(chimiddleware.RealIP)
	h.router.Use(middleware.NewZapMiddleware("router", l.Log))
	h.router.Use(middleware.Recoverer)

	return h
}

// GetRouter returns the router
func (h *Handler) GetRouter() chi.Router {
	return h.router
}
