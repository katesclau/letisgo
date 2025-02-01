package routes

import (
	"mnesis.com/frontend"
	"mnesis.com/pkg/server/endpoints"
)

// Routes is a map of API endpoints
func Get() *endpoints.APIRoutes {
	return &endpoints.APIRoutes{
		"/":        frontend.Index,
		"/health":  Health,
		"/profile": Profile,
	}
}
