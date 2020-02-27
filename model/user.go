package model

import (
	"log"

	"github.com/exlibris-fed/exlibris/key"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	registerModel(new(User))
}

// A User is a person interacting with the app. They may not be registered on this server.
type User struct {
	gorm.Model
	Username    string `gorm:"unique;not null;index"`
	DisplayName string `gorm:"not null"`
	Password    []byte `gorm:"not null" json:"-"`
	PrivateKey  []byte `gorm:"not null" json:"-"`
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
