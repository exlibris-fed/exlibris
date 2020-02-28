// Package model contains the models used by exlibris. Each model lives in its own file and should register itself via the registerModel function so that migrations will be applied.
package model

// A ContextKey is a key used to represent a model in a context
type ContextKey string
