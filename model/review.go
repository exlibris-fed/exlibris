package model

// Review models a book review
type Review struct {
	Base
	Book Book `gorm:"assocation_foreignkey:id"`
	User User
	Text string
}
