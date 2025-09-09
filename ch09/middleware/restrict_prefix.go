package middleware

import (
	"net/http"
	"path"
	"strings"
)

func RestrictPrefix(prefix string, next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			for p := range strings.SplitSeq(path.Clean(r.URL.Path), "/") {
				if strings.HasPrefix(p, prefix) {
					http.Error(w, "Not Found", http.StatusNotFound)
				}
			}

			next.ServeHTTP(w, r)
		},
	)
}
