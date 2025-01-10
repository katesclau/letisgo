package routes

import (
	"mnesis.com/pkg/server/endpoints"
)

// Routes is a map of API endpoints
func Get() *endpoints.APIRoutes {
	return &endpoints.APIRoutes{
		"/":       Root,
		"/health": Health,
	}
}
