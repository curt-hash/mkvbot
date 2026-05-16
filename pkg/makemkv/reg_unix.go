//go:build !windows

package makemkv

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var appKeyRe = regexp.MustCompile(`(?m)^app_Key\s*=.*$`)

// Reg registers key with MakeMKV by writing to ~/.MakeMKV/settings.conf.
func Reg(key string) error {
	dir := filepath.Join(os.Getenv("HOME"), ".MakeMKV")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir %q: %w", dir, err)
	}

	path := filepath.Join(dir, "settings.conf")
	line := fmt.Sprintf("app_Key = %q", key)

	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read %q: %w", path, err)
	}

	var out []byte
	if appKeyRe.Match(data) {
		out = appKeyRe.ReplaceAll(data, []byte(line))
	} else if len(data) > 0 {
		out = append(data, '\n')
		out = append(out, []byte(line+"\n")...)
	} else {
		out = []byte(line + "\n")
	}

	return os.WriteFile(path, out, 0644)
}
