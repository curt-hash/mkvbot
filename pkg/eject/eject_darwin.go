//go:build darwin

package eject

import (
	"context"
	"os/exec"
)

// Eject attempts to eject the disc.
func Eject(ctx context.Context, _ string) error {
	return exec.CommandContext(
		ctx,
		"drutil",
		"eject",
	).Run()
}
