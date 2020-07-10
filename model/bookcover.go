package model

import "github.com/google/uuid"

type BookCover struct {
	Base
	Book    Book
	BookID  string
	Cover   Cover
	CoverID uuid.UUID
}
