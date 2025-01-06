package makemkvcon

import "fmt"

var (
	ErrUnhandledLine = fmt.Errorf("unhandled line")
	ErrNotFound      = fmt.Errorf("not found")
)
