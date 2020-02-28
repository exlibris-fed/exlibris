package model

import (
	"crypto"
	"crypto/rsa"
	"fmt"
	"log"

	"github.com/exlibris-fed/exlibris/key"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	registerModel(new(User))
}

const (
	// ContextKeyRequestedUser  is the key to use for the username of the endpoint being requested.
	ContextKeyRequestedUser ContextKey = "username"

	// ContextKeyAuthenticatedUser is the key to use for a User object that was retrieved from a JWT. It should not be set until the JWT has been verified as being signed by the user specified in the `kid` field.
	ContextKeyAuthenticatedUser ContextKey = "authuser"

	// ContextKeyJWT is the key to use for a User's JWT in a context
	ContextKeyJWT ContextKey = "jwt"
)

// A User is a person interacting with the app. They may not be registered on this server.
type User struct {
	gorm.Model
	Username         string            `gorm:"unique;not null;index"`
	DisplayName      string            `gorm:"not null"`
	Password         []byte            `gorm:"not null" json:"-"`
	PrivateKey       []byte            `gorm:"not null" json:"-"`
	CryptoPrivateKey crypto.PrivateKey `gorm:"-"`
}

// SetPassword is used to hash the password the user wishes to use.
func (u *User) SetPassword(password string) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("error hashing password: " + err.Error())
		return
	}
	u.Password = hashed
}

// GenerateKeys is used on user registration to generate a private key for a user. It can theoretically be used to invalidate all existing tokens/sessions.
func (u *User) GenerateKeys() error {
	k, err := key.New()
	if err != nil {
		return err
	}
	bytes, err := key.SerializeRSAPrivateKey(k)
	if err != nil {
		return err
	}
	u.PrivateKey = bytes
	return nil
}

// IsPassword verifies that the specified password matches what's in the database.
func (u *User) IsPassword(password string) bool {
	return bcrypt.CompareHashAndPassword(u.Password, []byte(password)) == nil
}

func (u *User) ensureCryptoPrivateKey() {
	if u.CryptoPrivateKey == nil {
		pk, err := key.DeserializeRSAPrivateKey(u.PrivateKey)
		if err != nil {
			return
		}
		u.CryptoPrivateKey = pk
	}
}

// GenerateJWT generates a JWT for the user.
func (u *User) GenerateJWT() (string, error) {
	u.ensureCryptoPrivateKey()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"kid": u.Username, // TODO should we be generating urls as an ID?
	})
	t, err := token.SignedString(u.CryptoPrivateKey)
	if err != nil {
		return "", err
	}
	return t, nil
}

// ValidateJWT accepts a JWT and private key and verifies the token was signed by the key.
func (u *User) ValidateJWT(t string) bool {
	u.ensureCryptoPrivateKey()
	if u.CryptoPrivateKey == nil {
		// this may not be a user persisted in the database
		return false
	}

	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		pk, ok := u.CryptoPrivateKey.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("private key is %T, not *rsa.PrivateKey", u.CryptoPrivateKey)
		}
		return &pk.PublicKey, nil
	})
	if err != nil {
		log.Println("error validating JWT: " + err.Error())
		return false
	}
	return token.Valid
}
