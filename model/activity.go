package model

import "github.com/google/uuid"

// An APObject is an ActivityPub object (Read, Like, etc).
type APObject struct {
	Base
	User   User
	UserID uuid.UUID
	Read   Read
	ReadID uuid.UUID
}
