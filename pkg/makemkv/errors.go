package makemkv

import "fmt"

var (
	// ErrUnhandledLine is returned when attempting to parse a makemkvcon output
	// line with an unhandled prefix.
	ErrUnhandledLine = fmt.Errorf("unhandled line")

	// ErrNotFound is returned when something is not found.
	ErrNotFound = fmt.Errorf("not found")
)
