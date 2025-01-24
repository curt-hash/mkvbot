//go:build windows

package eject

import (
	"context"
	"fmt"
	"os/exec"
)

func Eject(ctx context.Context, volumeName string) error {
	return exec.CommandContext(
		ctx,
		"powershell",
		"-Command",
		fmt.Sprintf("(new-object -COM Shell.Application).NameSpace(17).ParseName('%s').InvokeVerb('Eject')", volumeName),
	).Run()
}
