package makemkv

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

// DriveScanLine represents a makemkvcon "DRV" output line, which describes a
// disc drive.
//
// See https://makemkv.com/developers/usage.txt.
type DriveScanLine struct {
	Index      int
	Visible    int
	Enabled    int
	Flags      int
	DriveName  string
	DiscTitle  string
	VolumeName string
}

var _ Line = (*DriveScanLine)(nil)

func (l *DriveScanLine) Kind() LineKind {
	return LineKindDriveScan
}

// ParseDriveScanLine parses the string that follows "DRV:" in the output of
// makemkvcon.
func ParseDriveScanLine(s string) (*DriveScanLine, error) {
	reader := csv.NewReader(strings.NewReader(s))
	reader.FieldsPerRecord = 7

	record, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("parse csv line %q: %w", s, err)
	}

	l := &DriveScanLine{
		DriveName:  record[4],
		DiscTitle:  record[5],
		VolumeName: record[6],
	}

	if l.Index, err = strconv.Atoi(record[0]); err != nil {
		return nil, fmt.Errorf("parse drive index %q: %w", record[0], err)
	}

	if l.Visible, err = strconv.Atoi(record[1]); err != nil {
		return nil, fmt.Errorf("parse field 1 value %q: %w", record[1], err)
	}

	if l.Enabled, err = strconv.Atoi(record[2]); err != nil {
		return nil, fmt.Errorf("parse field 2 value %q: %w", record[2], err)
	}

	if l.Flags, err = strconv.Atoi(record[3]); err != nil {
		return nil, fmt.Errorf("parse field 3 value %q: %w", record[3], err)
	}

	return l, nil
}
