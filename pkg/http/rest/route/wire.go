//+build wireinject

package route

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// ProvideHealthCheckRoute ...
func ProvideHealthCheckRoute(Gin *gin.Engine) HealthCheckRoute {
	return HealthCheckRoute{
		Gin: Gin,
	}
}

// ProvideAPIV1RouterGroup ...
func ProvideAPIV1RouterGroup(Gin *gin.Engine) APIV1RouterGroup {
	return APIV1RouterGroup{
		Gin: Gin,
	}
}

var Set = wire.NewSet(
	ProvideHealthCheckRoute,
	ProvideAPIV1RouterGroup,
	wire.Struct(new(Routes), "HealthCheckRoute", "APIV1RouterGroup"))

// InitializeRoutes
func InitializeRoutes(Gin *gin.Engine) Routes {
	panic(wire.Build(Set))
}
