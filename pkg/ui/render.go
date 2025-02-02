package ui

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/sirupsen/logrus"
	"mnesis.com/pkg/server/authorization"
	"mnesis.com/pkg/server/endpoints"
)

func Render(node func() templ.Component) endpoints.APIRouteEndpoint {
	return endpoints.APIRouteEndpoint{
		Handler: func(w http.ResponseWriter, r *http.Request) {
			logrus.Tracef("Processing path: %s", r.URL.Path)
			component := node()
			component.Render(r.Context(), w)
		},
		AuthorizationRole: authorization.None,
	}
}
