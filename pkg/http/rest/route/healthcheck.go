package route

import (
	"github.com/gin-gonic/gin"
	"github.com/marvincaspar/go-web-app-boilerplate/pkg/http/rest/handler"
)

// HealthCheckRoute ...
type HealthCheckRoute struct {
	Gin *gin.Engine
}

// Setup ...
func (r HealthCheckRoute) Setup() {
	r.Gin.GET("/health", handler.HealthCheckHandler)
}