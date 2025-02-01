//go:build linux

package eject

import (
	"context"
	"os/exec"
)

// Eject attempts to eject the disc identified by volumeName.
func Eject(ctx context.Context, volumeName string) error {
	return exec.CommandContext(
		ctx,
		"eject",
		volumeName,
	).Run()
}
