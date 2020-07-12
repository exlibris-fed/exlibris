package model

type Cover struct {
	Base
	Type   string `gorm:"not null"`
	URL    string `gorm:"unique;not null;index"`
	BookID string
}
