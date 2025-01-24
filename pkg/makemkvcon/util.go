package makemkvcon

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func ParseDuration(s string) (time.Duration, error) {
	tokens := strings.Split(s, ":")
	if n := len(tokens); n > 3 {
		return 0, fmt.Errorf("expected at most 3 tokens in duration string %q, got %d", s, n)
	}

	for len(tokens) < 3 {
		tokens = append([]string{"0"}, tokens...)
	}

	return time.ParseDuration(fmt.Sprintf("%sh%sm%ss", tokens[0], tokens[1], tokens[2]))
}
