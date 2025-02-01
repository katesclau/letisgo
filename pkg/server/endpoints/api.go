package endpoints

import (
	"net/http"

	"mnesis.com/pkg/server/authorization"
)

type APIDefinition struct {
	Name        string
	Description string
	Version     string
	Mux         *http.ServeMux
	Routes      *APIRoutes
}

type APIRouteEndpoint struct {
	Handler           http.HandlerFunc
	AuthorizationRole authorization.AuthorizationRole
}

type APIRoutes map[string]APIRouteEndpoint
