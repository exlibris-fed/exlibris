package model

import (
	"github.com/jinzhu/gorm"
)

// Review models a book review
type Review struct {
	gorm.Model
	FKBook int
	FKUser int
	Text   string
}
