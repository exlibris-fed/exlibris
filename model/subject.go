package model

// Subject for a book
type Subject struct {
	Base
	Books   []Book `gorm:"many2many:book_subjects"`
	Subject string `gorm:"unique"`
}

// NewSubject creates a new subject
func NewSubject(subject string) *Subject {
	return &Subject{
		Subject: subject,
	}
}
