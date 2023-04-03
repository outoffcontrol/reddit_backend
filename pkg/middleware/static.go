package middleware

import (
	"net/http"
	"net/url"
	"strings"
)

func StripUrlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "static") {
			next.ServeHTTP(w, r)
			return
		}
		p := "/"
		rp := "/"
		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = p
		r2.URL.RawPath = rp
		next.ServeHTTP(w, r2)
	})
}
