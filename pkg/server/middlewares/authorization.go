package middlewares

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"mnesis.com/pkg/server/authorization"
)

type AuthodizedRoutes = map[string]authorization.AuthorizationRole

type AuthorizationMiddlewareOptions struct {
	AuthorizedRoutes AuthodizedRoutes
}

type AuthorizationMiddleware struct {
	Options AuthorizationMiddlewareOptions
}

// AuthorizationMiddleware checks if the user is authorized
func (a AuthorizationMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
			"auth":   r.Header.Get("Authorization"),
		}).Info("AuthorizationMiddleware")

		if a.Options.AuthorizedRoutes[r.URL.Path] == authorization.None {
			next.ServeHTTP(w, r)
			return
		}

		// Placeholder for actual authorization logic
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
