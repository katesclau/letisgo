package ui

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/sirupsen/logrus"
	"mnesis.com/pkg/server/authorization"
	"mnesis.com/pkg/server/endpoints"
)

type Options struct {
	AuthorizationRole authorization.AuthorizationRole
	// Add more options here
}

type OptionFn func(*endpoints.APIRouteEndpoint) error

func RenderPage(node func() templ.Component, optFns ...OptionFn) endpoints.APIRouteEndpoint {
	c := endpoints.APIRouteEndpoint{}

	for _, optFn := range optFns {
		err := optFn(&c)
		if err != nil {
			return endpoints.APIRouteEndpoint{}
		}
	}
	c.Handler = func(w http.ResponseWriter, r *http.Request) {
		logrus.Tracef("[RenderPage] Processing path: %s", r.URL.Path)
		component := node()
		component.Render(r.Context(), w)
	}
	return c
}

func RenderComponent(handler func(w http.ResponseWriter, r *http.Request), optFns ...OptionFn) endpoints.APIRouteEndpoint {
	c := endpoints.APIRouteEndpoint{}

	for _, optFn := range optFns {
		err := optFn(&c)
		if err != nil {
			return endpoints.APIRouteEndpoint{}
		}
	}

	c.Handler = handler
	return c
}
