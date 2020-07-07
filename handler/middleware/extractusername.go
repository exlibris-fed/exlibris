package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/exlibris-fed/exlibris/model"

	"github.com/dgrijalva/jwt-go"
)

// ExtractUsername retrieves the authenticated username from a JWT. The expected header is:
//
//	Authorization: Bearer <JWT>
//
// where <JWT> is the token received during authentication. "Bearer" is case insensitive but is recommended to be title case.
//
// This does not require that authentication exist: it only retrieves it if it is present. To require it, use the Authenticated function.
func (m *Middleware) ExtractUsername(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		JWT := jwtFromAuth(r.Header.Get("Authorization"))
		if JWT == "" {
			next.ServeHTTP(w, r)
			return
		}

		token, err := jwt.Parse(JWT, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil || !token.Valid {
			next.ServeHTTP(w, r)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		username, ok := claims["kid"]
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		// "created" is also present in the claim. We could use this to expire a token after a certain amount of time, but I'm not sure how we should handle expiry/renewing.

		c := context.WithValue(r.Context(), model.ContextKeyAuthenticatedUsername, username)

		// r.Clone(c) was added in 1.13, but we're on 1.12
		next.ServeHTTP(w, r.WithContext(c))
	})
}
