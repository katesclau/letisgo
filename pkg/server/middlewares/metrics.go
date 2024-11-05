package middlewares

import (
	"log"
	"net/http"
	"time"
)

// MetricsMiddleware collects metrics for each request
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		// Placeholder for actual metrics collection logic
		log.Printf("Request to %s took %v", r.URL.Path, duration)
	})
}
