package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// HealthCheckHandler ...
func HealthCheckHandler(c *gin.Context) {
	// A very simple health check.

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	c.JSON(http.StatusOK, gin.H{
		"alive": true,
	})
}
