package model

import "github.com/google/uuid"

// Review models a book review
type Review struct {
	Base
	Book   Book
	BookID string
	User   User
	UserID uuid.UUID
	Text   string
}
