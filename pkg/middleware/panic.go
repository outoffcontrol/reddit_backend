package middleware

import (
	"net/http"
)

func (m *Middleware) Panic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				m.Logger.Z(r.Context()).Warn("recovered: ", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}
