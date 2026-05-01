package tracking

import "errors"

var (
	ErrActiveEntryExists = errors.New("active timer already exists")
	ErrNoActiveEntry     = errors.New("no active timer")
	ErrEmptyDescription  = errors.New("description cannot be empty")
)
