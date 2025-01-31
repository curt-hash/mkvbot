//go:build windows

package makemkvcon

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

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
