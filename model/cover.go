package model

type Cover struct {
	Base
	Type   string `gorm:"not null"`
	URL    string `gorm:"unique;not null;index"`
	Book   Book
	BookID string
}
