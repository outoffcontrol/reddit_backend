package middleware

import (
	"net/http"
	"time"
)

func (m *Middleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		m.Logger.Z(r.Context()).Infow(r.URL.Path,
			"method", r.Method,
			"remote_addr", r.RemoteAddr,
			"url", r.URL.Path,
			"work_time", time.Since(start),
		)
	})
}
