package makemkvcon

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

// From https://makemkv.com/developers/usage.txt:
//
// Current and total progress title
// PRGC:code,id,name
// PRGT:code,id,name
// code - unique message code
// id - operation sub-id
// name - name string
type progressLine struct {
	ID   int    `json:"id"`
	Code int    `json:"code"`
	Name string `json:"name"`
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
