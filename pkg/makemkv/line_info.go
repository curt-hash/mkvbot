package makemkv

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/curt-hash/mkvbot/pkg/makemkv/defs"
)

// InfoLine is the common representation of the "CINFO", "TINFO" and "SINFO"
// makemkvcon output lines, which describe an attribute of a disc, title, or
// stream.
//
// See https://makemkv.com/developers/usage.txt.
type InfoLine struct {
	prefix []int

	// ID is an integer that identifies the attribute.
	ID defs.Attr

	// Code is an integer that corresponds to Value, if Value is an enumeration.
	Code int

	// Value is the value of the attribute identified by ID.
	Value string
}

func (l *InfoLine) String() string {
	return fmt.Sprintf("%s: %s", l.ID, l.Value)
}

func parseInfoLine(numPrefixTokens int, s string) (*InfoLine, error) {
	reader := csv.NewReader(strings.NewReader(s))
	reader.FieldsPerRecord = numPrefixTokens + 3
	record, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("parse csv line %q: %w", s, err)
	}

	n := len(record) - 1
	ints := make([]int, 0, n)
	value, record := record[n], record[:n]
	for _, r := range record {
		i, err := strconv.Atoi(r)
		if err != nil {
			return nil, fmt.Errorf("parse %q: %w", r, err)
		}

		ints = append(ints, i)
	}

	l := &InfoLine{
		prefix: ints[:numPrefixTokens],
		ID:     defs.Attr(ints[n-2]),
		Code:   ints[n-1],
		Value:  value,
	}

	return l, nil
}

// Info is a slice of related InfoLines.
type Info []*InfoLine

// GetCode returns the Code of the InfoLine where ID matches id or ErrNotFound
// if such a line does not exist.
func (info Info) GetCode(id defs.Attr) (int, error) {
	for _, infoLine := range info {
		if infoLine.ID == id {
			return infoLine.Code, nil
		}
	}

	return 0, ErrNotFound
}

// GetCodeDefault returns the Code of the InfoLine where ID matches id or
// defaultValue if such a line does not exist.
func (info Info) GetCodeDefault(id defs.Attr, defaultValue int) int {
	v, err := info.GetCode(id)
	if err != nil {
		return defaultValue
	}

	return v
}

// GetAttr returns the Value of the InfoLine where ID matches id or ErrNotFound
// if such a line does not exist.
func (info Info) GetAttr(id defs.Attr) (string, error) {
	for _, infoLine := range info {
		if infoLine.ID == id {
			return infoLine.Value, nil
		}
	}

	return "", ErrNotFound
}

// GetAttrDefault returns the Value of the InfoLine where ID matches id or
// defaultValue if such a line does not exist.
func (info Info) GetAttrDefault(id defs.Attr, defaultValue string) string {
	v, err := info.GetAttr(id)
	if err != nil {
		return defaultValue
	}

	return v
}

// GetAttrInto is like GetAttr, except it also attempts to convert the Value to
// an integer.
func (info Info) GetAttrInt(id defs.Attr) (int, error) {
	v, err := info.GetAttr(id)
	if err != nil {
		return 0, err
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("parse %q: %w", v, err)
	}

	return n, nil
}

// GetAttrDuration is like GetAttr, except it also attempts to convert the
// value to a time.Duration.
func (info Info) GetAttrDuration(id defs.Attr) (time.Duration, error) {
	v, err := info.GetAttr(id)
	if err != nil {
		return 0, err
	}

	d, err := ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("parse %q: %w", v, err)
	}

	return d, nil
}
