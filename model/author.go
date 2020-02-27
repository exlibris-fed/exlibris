package model

import (
	"github.com/jinzhu/gorm"
)

func init() {
	registerModel(new(Author))
}

// An Author is someone who has written a Book.
type Author struct {
	gorm.Model
	Name string
}
