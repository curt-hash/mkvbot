package makemkv

import (
	"fmt"
	"strings"
	"time"
)

// ParseDuration parses a duration string with hours, minutes and seconds
// values separated by colons like "1:22:33".
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
