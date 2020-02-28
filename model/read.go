package model

import (
	"github.com/jinzhu/gorm"
)

func init() {
	registerModel(new(Read))
}

// Read is a many to many model describing a user who read a book
type Read struct {
	gorm.Model
	FKBook int
	FKUser int
}
