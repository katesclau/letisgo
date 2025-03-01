package endpoints

import (
	"net/http"

	"mnesis.com/pkg/server/authorization"
	"mnesis.com/pkg/server/session"
)

type APIDefinition struct {
	Name           string
	Description    string
	Version        string
	Mux            *http.ServeMux
	Routes         *APIRoutes
	SessionManager *session.SessionManager
}

type APIRouteEndpoint struct {
	Handler           http.HandlerFunc
	AuthorizationRole authorization.AuthorizationRole
}

type APIRoutes map[string]APIRouteEndpoint
