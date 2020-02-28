package model

import (
	"github.com/jinzhu/gorm"
)

// BookSubject for a book
type BookSubject struct {
	gorm.Model
	BookFK int
	SubjectFK int
}
