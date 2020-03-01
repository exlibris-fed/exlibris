package model

import (
	"time"
)

// An APObject is an ActivityPub object (Read, Like, etc).
type APObject struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	ID        string     `gorm:"primary key"`
	UserID    string     `gorm:"not null"`
	ReadID    uint       `gorm:"not null"`
}
