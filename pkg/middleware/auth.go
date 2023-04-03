package middleware

import (
	"context"

	"net/http"
	"reddit_backend/pkg/sessions"
	"reddit_backend/pkg/utils"
	"strings"
)

func (m *Middleware) Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		token = strings.Replace(token, "Bearer ", "", 1)
		user, err := m.Session.IsValid(r.Context(), token)
		if utils.HandleErr(w, r, err, 401, m.Logger) {
			return
		}
		ctx := context.WithValue(r.Context(), sessions.SessionKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
