package key

import (
	"crypto"
	"crypto/rsa"
	"fmt"
	"log"

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
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		pk, ok := k.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("private key is %T, not *rsa.PrivateKey", k)
		}
		return &pk.PublicKey, nil
	})
	if err != nil {
		log.Println("error validating JWT: " + err.Error())
		return false
	}
	return token.Valid
}
