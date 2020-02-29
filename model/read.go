package model

import (
	"github.com/jinzhu/gorm"
)

const (
	// ContextKeyRead is the context key to use for the read action
	ContextKeyRead ContextKey = "read"
)

// Read is a many to many model describing a user who read a book
type Read struct {
	gorm.Model
	FKBook string
	FKUser string
}
