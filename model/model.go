// Package model contains the models used by exlibris. Each model lives in its own file and should be initialized in main.go so that migrations run.
package model

import (
	"time"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/google/uuid"
)

// A ContextKey is a key used to represent a model in a context
type ContextKey string

// A Federater is a representation of a model as an ActivityPub object.
type Federater interface {
	ToType() vocab.Type
}

// Base set of attributes for a model
type Base struct {
	ID uuid.UUID `gorm:"primary_key"`
}

// BaseEvents is a set of attributes for events created, updated, deleted
type BaseEvents struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
