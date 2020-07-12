package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/exlibris-fed/exlibris/model"
	"github.com/google/uuid"
)

// WithUserModel takes the authenticated username from the context, if present, and populates the context with the User model. It does not require that it be present.
func (m *Middleware) WithUserModel(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()
		username, ok := c.Value(model.ContextKeyAuthenticatedUsername).(string)
		if !ok || username == "" {
			next.ServeHTTP(w, r)
			return
		}

		user, err := m.userRepo.GetByUsername(username)

		if err != nil || user.ID == uuid.Nil {
			log.Printf("user %s not present in database in UserModel middleware", username)
			next.ServeHTTP(w, r)
			return
		}

		c = context.WithValue(c, model.ContextKeyAuthenticatedUser, user)

		next.ServeHTTP(w, r.WithContext(c))
	})
}
