//go:build !windows

package eject

import (
	"context"
	"os/exec"
)

func Eject(ctx context.Context, volumeName string) error {
	return exec.CommandContext(
		ctx,
		"eject",
		volumeName,
	).Run()
}
