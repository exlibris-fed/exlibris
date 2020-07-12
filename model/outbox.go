package model

// An OutboxEntry represents an entry in a user's AP outbox.
type OutboxEntry struct {
	Base
	User       User
	UserID     string
	OutboxIRI  string
	URI        string
	Serialized string
}
