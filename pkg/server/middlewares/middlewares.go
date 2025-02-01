package middlewares

import (
	"net/http"
)

type MiddlewareOptions interface{}

type Interface interface {
	Handler(http.Handler) http.Handler
}

// ApplyMiddlewares applies all the middlewares to the given handler
func ApplyMiddlewares(handler http.Handler, middlewares ...Interface) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware.Handler(handler)
	}
	return handler
}
