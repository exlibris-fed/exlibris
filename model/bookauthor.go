package model

// BookAuthor is a many to many model describing the authors for books
type BookAuthor struct {
	Base
	Book   Book
	Author Author
}
