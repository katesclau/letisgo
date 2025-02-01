package routes

import (
	"mnesis.com/frontend"
	"mnesis.com/pkg/server/endpoints"
	"mnesis.com/pkg/ui"
)

// Routes is a map of API endpoints
func Get() *endpoints.APIRoutes {
	return &endpoints.APIRoutes{
		"/":        ui.GetRenderAPIRouteEndpoint(frontend.Index),
		"/health":  Health,
		"/profile": Profile,
	}
}
