package server

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := newStatusResponseWriter(w)
		defer func() {
			duration := time.Since(start)
			log.Printf(
				"[%s] %s %s %d %s",
				r.Method,
				r.RequestURI,
				r.RemoteAddr,
				wrapped.status,
				duration,
			)
		}()
		next.ServeHTTP(wrapped, r)
	})
}


type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func newStatusResponseWriter(w http.ResponseWriter) *statusResponseWriter {
	return &statusResponseWriter{
		ResponseWriter: w,
		status:         http.StatusOK,
	}
}

func (w *statusResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}