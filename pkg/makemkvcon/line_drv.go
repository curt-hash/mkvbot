package makemkvcon

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

// From https://makemkv.com/developers/usage.txt:
//
// DRV:index,visible,enabled,flags,drive name,disc name
// index - drive index
// visible - set to 1 if drive is present
// enabled - set to 1 if drive is accessible
// flags - media flags, see AP_DskFsFlagXXX in apdefs.h
// drive name - drive name string
// disc name - disc name string
//
// Errata:
//   - Field 1 (visible?) is not boolean (e.g., 0, 1, 2)
//   - Field 2 (enabled?) is not boolean (e.g., 999)
//   - Undocumented title field between drive and disc name
type DriveScanLine struct {
	Index      int    `json:"index"`
	Field1     int    `json:"visible"`
	Field2     int    `json:"enabled"`
	Field3     int    `json:"flags"`
	DriveName  string `json:"drive_name"`
	DiscTitle  string `json:"disc_title"`
	VolumeName string `json:"volume_name"`
}

var _ Line = (*DriveScanLine)(nil)

func (l *DriveScanLine) Kind() LineKind {
	return LineKindDriveScan
}

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

	if l.Field1, err = strconv.Atoi(record[1]); err != nil {
		return nil, fmt.Errorf("parse field 1 value %q: %w", record[1], err)
	}

	if l.Field2, err = strconv.Atoi(record[2]); err != nil {
		return nil, fmt.Errorf("parse field 2 value %q: %w", record[2], err)
	}

	if l.Field3, err = strconv.Atoi(record[3]); err != nil {
		return nil, fmt.Errorf("parse field 3 value %q: %w", record[3], err)
	}

	return l, nil
}
