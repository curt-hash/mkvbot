package makemkv

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"os"
	"strings"
)

// ParseLine parses a single line of output from makemkvcon.
func ParseLine(s string) (*Line, error) {
	return lineParser.ParseString("", strings.TrimSpace(s))
}

// ParseLines parses makemkvcon output lines from r. It returns a sequence of
// [*Line, error] where either Line is a parsed line or err is non-nil. The
// sequence ends after all lines have been parsed and r returns EOF. Individual
// line parsing errors do not trigger an early return.
func ParseLines(r io.Reader) iter.Seq2[*Line, error] {
	return func(yield func(*Line, error) bool) {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			s := scanner.Text()
			line, err := ParseLine(s)
			if err != nil {
				err = fmt.Errorf("parse line %q: %w", s, err)
			}
			if !yield(line, err) {
				return
			}
		}

		if err := scanner.Err(); err != nil {
			yield(nil, fmt.Errorf("scan lines from stdout: %w", err))
		}
	}
}

// LineIterator is a generic type that represents the lines output by a
// makemkvcon command and the generic final result.
type LineIterator[T any] struct {
	Seq    iter.Seq2[*Line, error]
	result T
	err    error
}

// GetResult returns the final result of the command.
func (li *LineIterator[T]) GetResult() (T, error) {
	return li.result, li.err
}

// ParseOutput parses multi-line output from makemkvcon.
func ParseOutput(s string) (*Output, error) {
	return outputParser.ParseString("", s)
}

// ParseFile parses the multi-line output of makemkvcon from a file.
func ParseFile(path string) (*Output, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", path, err)
	}
	defer f.Close()

	return outputParser.Parse(path, f)
}

// Grammar returns an EBNF representation of the supported grammar.
func Grammar() string {
	return outputParser.String()
}
