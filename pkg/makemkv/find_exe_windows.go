//go:build windows

package makemkv

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// FindExe attempts to return the path of the makemkvcon executable on Windows
// operating systems.
func FindExe() (string, error) {
	exe := "makemkvcon64.exe"
	if path, err := exec.LookPath(exe); err == nil {
		return path, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}

	volume := filepath.VolumeName(wd)
	return filepath.Join(volume+"\\", "Program Files (x86)", "MakeMKV", exe), nil
}
