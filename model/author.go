package model

import (
	"github.com/jinzhu/gorm"
)

// An Author is someone who has written a Book.
type Author struct {
	gorm.Model
	Name string `json:"name"`
}
