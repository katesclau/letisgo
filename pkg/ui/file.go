package ui

import (
	"net/http"
	"path/filepath"

	"mnesis.com/pkg/server/authorization"
	"mnesis.com/pkg/server/endpoints"
)

var File = endpoints.RouteEndpoint{
	Handler: http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			filePath := r.URL.Path[len("/static/"):]
			fullPath := filepath.Join(".", "frontend/static", filePath)
			http.ServeFile(w, r, fullPath)
		}),
	AuthorizationRole: authorization.None,
}
