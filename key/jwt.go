package key

import (
	"crypto"

	"github.com/dgrijalva/jwt-go"
)

// GenerateJWT generates a JWT using the specified key.
func GenerateJWT(k crypto.PrivateKey) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	t, err := token.SignedString(k)
	if err != nil {
		return "", err
	}
	return t, nil
}

// ValidateJWT accepts a JWT and private key and verifies the token was signed by the key.
func ValidateJWT(t string, k crypto.PrivateKey) bool {
	// TODO
	// https://godoc.org/github.com/dgrijalva/jwt-go#example-Parse--Hmac
	return false
}
