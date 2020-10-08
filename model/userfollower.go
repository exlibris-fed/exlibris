package model

import (
	"github.com/google/uuid"
)

// A Follower is the IRI of someone who follows a user.
type Follower struct {
	ID     string `gorm:"primary_key"`
	User   User   `gorm:"association_autoupdate:false"`
	UserID uuid.UUID
}
