package model

// An OutboxEntry represents an entry in a user's AP outbox.
type OutboxEntry struct {
	Base
	User       User
	Serialized string `gorm:"primary_key"`
}
