package model

// A Follower is the IRI of someone who follows a user.
type Follower struct {
	ID     string `gorm:"primary_key"`
	User   User
	UserID string
}
