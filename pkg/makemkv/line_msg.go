package makemkv

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

// MessageLine represents a makemkvcon "MSG" output line, which is an
// informational logging line.
//
// See https://makemkv.com/developers/usage.txt.
type MessageLine struct {
	Code      int
	Flags     int
	NumParams int
	Message   string
	Format    string
	Params    []string
}

var _ Line = (*MessageLine)(nil)

func (l *MessageLine) Kind() LineKind {
	return LineKindMessage
}

// ParseMessageLine parses the string that follows "MSG:" in the output of
// makemkvcon.
func ParseMessageLine(s string) (*MessageLine, error) {
	tokens := strings.SplitN(s, ",", 4)
	if len(tokens) != 4 {
		return nil, fmt.Errorf("expected 4 tokens, got %d", len(tokens))
	}

	numParams, err := strconv.Atoi(tokens[2])
	if err != nil {
		return nil, fmt.Errorf("parse count %q: %w", tokens[2], err)
	}

	reader := csv.NewReader(strings.NewReader(s))
	reader.FieldsPerRecord = 5 + numParams

	record, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("parse csv line %q: %w", s, err)
	}

	l := &MessageLine{
		NumParams: numParams,
		Message:   record[3],
		Format:    record[4],
		Params:    record[5:],
	}

	if l.Code, err = strconv.Atoi(record[0]); err != nil {
		return nil, fmt.Errorf("parse code %q: %w", record[0], err)
	}

	if l.Flags, err = strconv.Atoi(record[1]); err != nil {
		return nil, fmt.Errorf("parse flags %q: %w", record[1], err)
	}

	return l, nil
}
