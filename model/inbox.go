package model

import "github.com/google/uuid"

// An InboxEntry represents an entry in a user's AP inbox.
type InboxEntry struct {
	Base
	User       User `gorm:"association_autoupdate:false"`
	UserID     uuid.UUID
	InboxIRI   string
	URI        string
	Serialized string
}
