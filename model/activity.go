package model

// An APObject is an ActivityPub object (Read, Like, etc).
type APObject struct {
	Base
	User User
	Read Read
}
