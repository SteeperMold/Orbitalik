package middleware

import (
	"net/http"
	"time"

	applog "github.com/SteeperMold/Orbitalik/common/go/log"
)

type logRecord struct {
	statusCode int
	http.ResponseWriter
}

func (lr *logRecord) WriteHeader(code int) {
	lr.statusCode = code
	lr.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware returns an HTTP middleware that logs incoming requests and responses.
// It generates or propagates a correlation ID (X-Correlation-ID) for tracing requests,
// logs the request start with method and URI, then logs the response status and duration.
func LoggingMiddleware(logger applog.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cid := r.Header.Get("X-Correlation-ID")

			start := time.Now()
			lr := &logRecord{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}
			next.ServeHTTP(lr, r)

			duration := time.Since(start)
			logger.Info("http request",
				applog.NewField("correlation_id", cid),
				applog.NewField("method", r.Method),
				applog.NewField("uri", r.RequestURI),
				applog.NewField("status", lr.statusCode),
				applog.NewField("duration", duration),
			)
		})
	}
}
