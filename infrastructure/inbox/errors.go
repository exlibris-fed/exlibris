package inbox

import (
	"errors"
)

var (
	// ErrNotFound is returned when a record cannot be found.
	ErrNotFound = errors.New("inbox could not be found")
	// ErrEntryNotCreated is returned when a record cannot be created.
	ErrEntryNotCreated = errors.New("inbox entry could not be saved")
)
