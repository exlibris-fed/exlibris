package model

import (
	"github.com/jinzhu/gorm"
)

// Read is a many to many model describing a user who read a book
type Read struct {
	gorm.Model
	FKBook int
	FKUser int
}
