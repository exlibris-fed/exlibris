package model

import (
	"github.com/jinzhu/gorm"
)

func init() {
	registerModel(new(BookAuthor))
}

// BookAuthor is a many to many model describing the authors for books
type BookAuthor struct {
	gorm.Model
	FKBook int
	FKAuthor int
}
