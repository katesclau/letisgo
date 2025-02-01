package ui

import (
	"net/http"

	"github.com/a-h/templ"
	"mnesis.com/pkg/server/authorization"
	"mnesis.com/pkg/server/endpoints"
)

func GetRenderAPIRouteEndpoint(node func() templ.Component) endpoints.APIRouteEndpoint {
	return endpoints.APIRouteEndpoint{
		Handler: func(w http.ResponseWriter, r *http.Request) {
			component := node()
			component.Render(r.Context(), w)
		},
		AuthorizationRole: authorization.None,
	}
}
