package model

import (
	"github.com/jinzhu/gorm"
)

// A Book is something that can be read. Currently this only supports things which are in the Library of Congress API, but eventually it'd be great to support fanfiction and other online-only sources.
type Book struct {
	gorm.Model
	Key       string   `gorm:"PRIMARY_KEY" json:"key"`
	Title     string   `gorm:"not null;index" json:"title"`
	Authors   []Author `json:"authors"`
	Published int      `json:"published"`
	ISBN      string   `json:"isbn,omitempty"`
	Subjects  []string `json:"subjects"` // TODO do we want to index this?
}
