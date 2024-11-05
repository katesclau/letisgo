package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// TracingMiddleware adds a unique trace ID to each request
func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.New().String()
		ctx := context.WithValue(r.Context(), "TraceID", traceID)
		r = r.WithContext(ctx)
		w.Header().Set("X-Trace-ID", traceID)
		next.ServeHTTP(w, r)
	})
}
