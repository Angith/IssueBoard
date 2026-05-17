package middleware

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/sirupsen/logrus"
)

// responseWriter captures the status code for logging
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggerMiddleware logs incoming HTTP requests and their duration/status
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		logrus.WithFields(logrus.Fields{
			"method":   r.Method,
			"path":     r.URL.Path,
			"status":   rw.status,
			"duration": duration,
			"ip":       r.RemoteAddr,
		}).Info("HTTP Request")
	})
}

// RecoveryMiddleware catches panics, logs them, and returns a 500 error instead of crashing
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logrus.WithFields(logrus.Fields{
					"method": r.Method,
					"path":   r.URL.Path,
					"error":  err,
				}).Error("Panic recovered during HTTP request")
				
				// Optional: Log the stack trace in debug mode
				logrus.Debugf("Stack trace: %s", debug.Stack())

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
