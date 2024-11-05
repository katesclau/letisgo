package middlewares

import (
	"net/http"
)

// ApplyMiddlewares applies all the middlewares to the given handler
func ApplyMiddlewares(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
