package model

import (
	"crypto"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/exlibris-fed/exlibris/key"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	// ContextKeyRequestedUser  is the key to use for the username of the endpoint being requested.
	ContextKeyRequestedUser ContextKey = "username"

	// ContextKeyAuthenticatedUsername is the key to use for a username that was retrieved and validated from the Authorization header. It should not be set until the header has been verified.
	ContextKeyAuthenticatedUsername ContextKey = "authusername"

	// ContextKeyAuthenticatedUser is the key to use for a User object that was retrieved from a JWT. It should not be set until the JWT has been verified as being signed by the user specified in the `kid` field.
	ContextKeyAuthenticatedUser ContextKey = "authuser"

	// ContextKeyJWT is the key to use for a User's JWT in a context
	ContextKeyJWT ContextKey = "jwt"
)

// A User is a person interacting with the app. They may not be registered on this server.
type User struct {
	Base
	HumanID          string `gorm:"unique;not null;index"`
	Username         string `gorm:"unique;not null;index"`
	DisplayName      string `gorm:"not null"`
	Email            string `gorm:"not null"`
	Password         []byte `json:"-"`
	PrivateKey       []byte `json:"-"`
	Summary          string
	CryptoPrivateKey crypto.PrivateKey `gorm:"-"`
	Local            bool              `json:"-"`
	Verified         bool              `json:"-"`
}

// NewUser creates a user and handles generating the ID, key and hashed password.
func NewUser(username, password, email, displayName string) (*User, error) {
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		return nil, fmt.Errorf("DOMAIN env variable not set")
	}

	u := User{
		Base: Base{
			ID: uuid.New(),
		},
		HumanID:     fmt.Sprintf("%s/@%s", domain, username),
		Username:    username,
		Email:       email,
		DisplayName: displayName,
	}
	u.SetPassword(password)
	if err := u.GenerateKeys(); err != nil {
		return nil, err
	}
	return &u, nil
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

// IRI returns a url representing the user's profile
func (u *User) IRI() *url.URL {
	URL, err := url.Parse(fmt.Sprintf("https://%s", u.HumanID))
	if err != nil {
		log.Printf("error creating IRI for user %s (%s): %s", u.ID, u.Username, err)
		return nil
	}
	return URL
}

// OutboxIRI returns a url representing the user's outbox
func (u *User) OutboxIRI() *url.URL {
	URL, err := url.Parse(fmt.Sprintf("https://%s/outbox", u.HumanID))
	if err != nil {
		log.Printf("error creating outbox IRI for user %s (%s): %s", u.ID, u.Username, err)
		return nil
	}
	return URL
}

// InboxIRI returns a url representing the user's inbox
func (u *User) InboxIRI() *url.URL {
	URL, err := url.Parse(fmt.Sprintf("https://%s/inbox", u.HumanID))
	if err != nil {
		log.Printf("error creating inbox IRI for user %s (%s): %s", u.HumanID, u.Username, err)
		return nil
	}
	return URL
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"kid":     u.Username, // TODO should we be generating urls as an ID?
		"created": time.Now(), //.Unix,
	})
	t, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}
	return t, nil
}

// ToType returns a representation of a user as an ActivityPub object.
func (u *User) ToType() vocab.Type {
	user := streams.NewActivityStreamsPerson()

	URL, err := url.Parse(u.HumanID)
	if err == nil {
		id := streams.NewJSONLDIdProperty()
		id.SetIRI(URL)
		user.SetJSONLDId(id)
	}

	name := streams.NewActivityStreamsNameProperty()
	name.AppendXMLSchemaString(u.DisplayName)
	user.SetActivityStreamsName(name)

	username := streams.NewActivityStreamsPreferredUsernameProperty()
	username.SetXMLSchemaString(u.Username)
	user.SetActivityStreamsPreferredUsername(username)

	return user
}
