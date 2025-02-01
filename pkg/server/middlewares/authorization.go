package middlewares

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"
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
		ctx := r.Context()
		log.WithFields(log.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
			"auth":   r.Header.Get("Authorization"),
		}).Trace("AuthorizationMiddleware")

		if r.Header.Get("Authorization") == "" {
			ctx = context.WithValue(ctx, "user", "none")
			log.WithFields(log.Fields{
				"user": ctx.Value("user"),
			}).Trace("[AuthorizationMiddleware] Defined user")
		}

		if a.Options.AuthorizedRoutes[r.URL.Path] == authorization.None {
			next.ServeHTTP(w, r.Clone(ctx))
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
