package model

import "github.com/google/uuid"

// An InboxEntry represents an entry in a user's AP inbox.
type InboxEntry struct {
	Base
	User       User
	UserID     uuid.UUID
	InboxIRI   string
	URI        string
	Serialized string
}
