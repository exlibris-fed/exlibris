// Package model contains the models used by exlibris. Each model lives in its own file and should be initialized in main.go so that migrations run.
package model

import (
	"github.com/go-fed/activity/streams/vocab"
)

// A ContextKey is a key used to represent a model in a context
type ContextKey string

// A Federater is a representation of a model as an ActivityPub object.
type Federater interface {
	ToType() vocab.Type
}
