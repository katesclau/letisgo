package middlewares

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type MetricsMiddleware struct{}

// MetricsMiddleware collects metrics for each request
func (m MetricsMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		// Placeholder for actual metrics collection logic
		logrus.Infof("Request to %s took %v", r.URL.Path, duration)
	})
}
