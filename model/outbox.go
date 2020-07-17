package model

import "github.com/google/uuid"

// An OutboxEntry represents an entry in a user's AP outbox.
type OutboxEntry struct {
	Base
	User       User `gorm:"association_autoupdate:false"`
	UserID     uuid.UUID
	Serialized string `gorm:"primary_key"`
}
