//go:build windows

package makemkv

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

// Reg registers key with MakeMKV by writing directly to the Windows registry.
func Reg(key string) error {
	k, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\MakeMKV`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("open HKCU\\Software\\MakeMKV: %w", err)
	}
	defer k.Close()
	if err := k.SetStringValue("app_Key", key); err != nil {
		return fmt.Errorf("set app_Key: %w", err)
	}
	return nil
}
