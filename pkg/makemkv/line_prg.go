package makemkv

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

// progressLine is the common representation for makemkvcon "PRGC" and "PRGT"
// output lines, which describe the current and overall task, respectively.
type progressLine struct {
	ID   int
	Code int

	// Name is the name of the task.
	Name string
}

func parseProgressLine(s string) (*progressLine, error) {
	reader := csv.NewReader(strings.NewReader(s))
	reader.FieldsPerRecord = 3
	record, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("parse csv line %q: %w", s, err)
	}

	ints := make([]int, 0, 2)
	for _, r := range record[:2] {
		i, err := strconv.Atoi(r)
		if err != nil {
			return nil, fmt.Errorf("parse %q: %w", r, err)
		}

		ints = append(ints, i)
	}

	return &progressLine{
		ID:   ints[0],
		Code: ints[1],
		Name: record[2],
	}, nil
}
