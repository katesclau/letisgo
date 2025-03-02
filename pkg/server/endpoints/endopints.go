package endpoints

import (
	"net/http"

	"mnesis.com/pkg/server/authorization"
)

type Endpoints struct {
	Name        string
	Description string
	Version     string
	Mux         *http.ServeMux
	Routes      *Routes
}

type RouteEndpoint struct {
	Handler           http.HandlerFunc
	AuthorizationRole authorization.AuthorizationRole
}

type Routes map[string]RouteEndpoint

func GetAuthorizationRoutes(routes *Routes) map[string]authorization.AuthorizationRole {
	authorizationRoutes := make(map[string]authorization.AuthorizationRole)
	for path, endpoint := range *routes {
		authorizationRoutes[path] = endpoint.AuthorizationRole
	}
	return authorizationRoutes
}
