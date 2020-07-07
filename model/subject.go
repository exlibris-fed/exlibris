package model

import "github.com/google/uuid"

// Subject for a book
type Subject struct {
	Base
	Books          []Book `gorm:"many2many:book_subjects"`
	BookSubjectsID uuid.UUID
	Subject        string
}
