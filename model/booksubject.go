package model

// BookSubject for a book
type BookSubject struct {
	Base
	Book      Book
	BookID    string
	Subject   Subject
	SubjectID string
}
