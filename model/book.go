package model

import (
	"github.com/jinzhu/gorm"
)

func init() {
	registerModel(new(Book))
}

// A Book is something that can be read. Currently this only supports things which are in the Library of Congress API, but eventually it'd be great to support fanfiction and other online-only sources.
type Book struct {
	gorm.Model
	Title   string `gorm:"not null;index"`
	Authors []Author
}
