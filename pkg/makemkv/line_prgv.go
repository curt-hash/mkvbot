package makemkv

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

// ProgressBarLine represents a "PRGV" output line, which describes the
// progress of a task and sub-task.
//
// See https://makemkv.com/developers/usage.txt.
type ProgressBarLine struct {
	// Current represents the progress of the current sub-task.
	Current int

	// Total represents the progress of the overall task.
	Total int

	// Max is a constant denominator used to calculate the progress percentage.
	Max int
}

func (l *ProgressBarLine) Kind() LineKind {
	return LineKindProgressBar
}

// CurrentProgress returns the progress of the current sub-task as a
// percentage.
func (l *ProgressBarLine) CurrentProgress() float64 {
	return float64(l.Current) / float64(l.Max)
}

// TotalProgress returns the progress of the overall task as a percentage.
func (l *ProgressBarLine) TotalProgress() float64 {
	return float64(l.Total) / float64(l.Max)
}

// ParseProgressBarLine parses the string that follows "PRGV:" in the output of
// makemkvcon.
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
