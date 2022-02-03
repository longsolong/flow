package route

// Routes contains multiple routes
type Routes struct {
	HealthCheckRoute
	APIV1RouterGroup
}

// Route interface
type Route interface {
	Setup()
}

// Setup all the route
func (r Routes) Setup() {
	r.HealthCheckRoute.Setup()
	r.APIV1RouterGroup.Setup()
}

