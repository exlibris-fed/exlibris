package model

import (
	"github.com/jinzhu/gorm"
)

// Read is a many to many model describing a user who read a book
type Read struct {
	gorm.Model
	FKBook string`gorm:"primary_key;not null"`
	FKUser string`gorm:"primary_key;not null"`
}
