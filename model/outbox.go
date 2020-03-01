package model

import (
	"time"
)

// An OutboxEntry represents an entry in a user's AP outbox.
type OutboxEntry struct {
	CreatedAt  time.Time
	UserID     string `gorm:"primary_key"`
	Serialized string `gorm:"primary_key"`
}
