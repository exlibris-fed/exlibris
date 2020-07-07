package middleware

import (
	"net/http"

	"github.com/exlibris-fed/exlibris/model"
)

// Authenticated requires that the username added in ExtractUsername be present.
func (m *Middleware) Authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()
		if _, ok := c.Value(model.ContextKeyAuthenticatedUsername).(string); !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
