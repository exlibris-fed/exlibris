package model

import (
	"github.com/jinzhu/gorm"
)

func init() {
	registerModel(new(Review))
}

// Review models a book review
type Review struct {
	gorm.Model
	FKBook int
	FKUser int
	Text string
}
