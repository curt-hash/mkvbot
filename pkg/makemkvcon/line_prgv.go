package makemkvcon

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

// From https://makemkv.com/developers/usage.txt:
//
// Progress bar values for current and total progress
// PRGV:current,total,max
// current - current progress value
// total - total progress value
// max - maximum possible value for a progress bar, constant
type ProgressBarLine struct {
	Current int `json:"current"`
	Total   int `json:"total"`
	Max     int `json:"max"`
}

func (l *ProgressBarLine) Kind() LineKind {
	return LineKindProgressBar
}

func (l *ProgressBarLine) CurrentProgress() float64 {
	return float64(l.Current) / float64(l.Max)
}

func (l *ProgressBarLine) TotalProgress() float64 {
	return float64(l.Total) / float64(l.Max)
}

func ParseProgressBarLine(s string) (*ProgressBarLine, error) {
	reader := csv.NewReader(strings.NewReader(s))
	reader.FieldsPerRecord = 3
	record, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("parse csv line %q: %w", s, err)
	}

	ints := make([]int, 0, 3)
	for _, r := range record {
		i, err := strconv.Atoi(r)
		if err != nil {
			return nil, fmt.Errorf("parse %q: %w", r, err)
		}

		ints = append(ints, i)
	}

	return &ProgressBarLine{
		Current: ints[0],
		Total:   ints[1],
		Max:     ints[2],
	}, nil
}
