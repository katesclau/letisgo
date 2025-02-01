package middlewares

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type LoggingMiddleware struct{}

func (lm LoggingMiddleware) Options() MiddlewareOptions {
	return nil
}

// LoggingMiddleware logs the details of each request
func (lm LoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logrus.Infof("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		logrus.Infof("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}
