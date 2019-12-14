package middleware

import (
	"fmt"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap/zapcore"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type zaplogger struct {
	logZ *zap.Logger
	name string
}

// NewZapMiddleware returns a new Zap Middleware handler.
func NewZapMiddleware(name string, logger *zap.Logger) func(next http.Handler) http.Handler {
	return zaplogger{
		logZ: logger,
		name: name,
	}.middleware
}

func (c zaplogger) middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var requestID string
		if reqID := r.Context().Value(middleware.RequestIDKey); reqID != nil {
			requestID = reqID.(string)
		}
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		latency := time.Since(start)

		if c.logZ != nil {
			fields := []zapcore.Field{
				zap.Int("status", ww.Status()),
				zap.Duration("took", latency),
				zap.Int64(fmt.Sprintf("measure#%s.latency", c.name), latency.Nanoseconds()),
				zap.String("remote", r.RemoteAddr),
				zap.String("request", r.RequestURI),
				zap.String("method", r.Method),
			}
			if requestID != "" {
				fields = append(fields, zap.String("request-id", requestID))
			}
			c.logZ.Info("request completed", fields...)
		}
	}
	return http.HandlerFunc(fn)
}