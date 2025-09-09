package middleware

import (
	"net/http"
	"path"
	"strings"
)

// AllowPrefix is the inverse of the [RestrictPrefix] middleware. It permits
// requests for only an allowed list of resources.
func AllowPrefix(prefix string, next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			for p := range strings.SplitSeq(path.Clean(r.URL.Path), "/") {
				if strings.HasPrefix(p, prefix) {
					next.ServeHTTP(w, r)
				}
			}

			http.Error(w, "Not Found", http.StatusNotFound)
		},
	)
}
