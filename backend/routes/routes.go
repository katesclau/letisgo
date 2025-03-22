package routes

import (
	"mnesis.com/frontend/pages"
	"mnesis.com/pkg/server/endpoints"
	"mnesis.com/pkg/ui"
)

// Routes is a map of API endpoints
func Get() *endpoints.Routes {
	return &endpoints.Routes{
		// Pages
		"GET /":        ui.RenderPage(pages.Home),
		"GET /about":   ui.RenderPage(pages.About),
		"GET /contact": ui.RenderPage(pages.Contact),
		"GET /static/": ui.File,

		// Health
		"GET /health": Health,

		// User
		"POST /user/login":    Login,
		"POST /user/register": Register,
		"POST /user/forgot":   Forgot,
	}
}
