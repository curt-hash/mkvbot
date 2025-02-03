//go:build linux || darwin

package makemkv

import (
	"os/exec"
)

// FindExe attempts to return the path of the makemkvcon executable on Linux
// and Darwin operating systems.
func FindExe() (string, error) {
	return exec.LookPath("makemkvcon")
}
