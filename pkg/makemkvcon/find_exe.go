//go:build linux || darwin

package makemkvcon

import (
	"os/exec"
)

func find_exe() (string, error) {
	exe := "makemkvcon"
	if path, err := exec.LookPath(exe); err == nil {
		return path, nil
	}

	return exe, nil
}
