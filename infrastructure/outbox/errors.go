package outbox

import (
	"errors"
)

var (
	// ErrNotFound is returned when a record cannot be found.
	ErrNotFound = errors.New("outbox could not be found")
	// ErrEntryNotCreated is returned when a record cannot be created.
	ErrEntryNotCreated = errors.New("outbox entry could not be saved")
)
