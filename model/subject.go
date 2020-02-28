package model

import (
	"github.com/jinzhu/gorm"
)

// Subject for a book
type Subject struct {
	gorm.Model
	Subject string
}
