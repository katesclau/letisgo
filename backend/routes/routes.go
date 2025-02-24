package routes

import (
	"mnesis.com/frontend/pages"
	"mnesis.com/pkg/server/endpoints"
	"mnesis.com/pkg/ui"
)

// Routes is a map of API endpoints
func Get() *endpoints.APIRoutes {
	return &endpoints.APIRoutes{
		"GET /health":  Health,
		"POST /login":  Login,
		"GET /static/": ui.File,
		"GET /":        ui.RenderPage(pages.Home),
	}
}
