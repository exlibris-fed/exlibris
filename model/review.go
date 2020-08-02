package model

import "github.com/google/uuid"

// Review models a book review
type Review struct {
	Base
	Book   Book      `gorm:"association_autoupdate:false"`
	BookID string    `gorm:"index"`
	User   User      `gorm:"association_autoupdate:false"`
	UserID uuid.UUID `gorm:"index"`
	Text   string
}
