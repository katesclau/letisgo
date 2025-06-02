package ui

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/katesclau/letisgo/internal/server/authorization"
	"github.com/katesclau/letisgo/internal/server/endpoints"
	"github.com/sirupsen/logrus"
)

type Options struct {
	AuthorizationRole authorization.AuthorizationRole
	// Add more options here
}

type OptionFn func(*endpoints.RouteEndpoint) error

func RenderPage(node func() templ.Component, optFns ...OptionFn) endpoints.RouteEndpoint {
	c := endpoints.RouteEndpoint{}

	for _, optFn := range optFns {
		err := optFn(&c)
		if err != nil {
			return endpoints.RouteEndpoint{}
		}
	}
	c.Handler = func(w http.ResponseWriter, r *http.Request) {
		logrus.Tracef("[RenderPage] Processing path: %s", r.URL.Path)
		component := node()
		component.Render(r.Context(), w)
	}
	return c
}

func RenderComponent(handler func(w http.ResponseWriter, r *http.Request), optFns ...OptionFn) endpoints.RouteEndpoint {
	c := endpoints.RouteEndpoint{}

	for _, optFn := range optFns {
		err := optFn(&c)
		if err != nil {
			return endpoints.RouteEndpoint{}
		}
	}

	c.Handler = handler
	return c
}
