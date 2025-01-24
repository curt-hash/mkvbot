package makemkvcon

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/curt-hash/mkvbot/pkg/makemkvcon/defs"
)

// From https://makemkv.com/developers/usage.txt:
//
// CINFO:id,code,value
// TINFO:id,code,value
// SINFO:id,code,value
// id - attribute id, see AP_ItemAttributeId in apdefs.h
// code - message code if attribute value is a constant string
// value - attribute value
type InfoLine struct {
	prefix []int
	ID     defs.Attr `json:"id"`
	Code   int       `json:"code"`
	Value  string    `json:"value"`
}

func (l *InfoLine) String() string {
	return fmt.Sprintf("%s: %s", defs.Attr(l.ID), l.Value)
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

type Info []*InfoLine

func (info Info) GetAttr(id defs.Attr) (string, error) {
	for _, infoLine := range info {
		if infoLine.ID == id {
			return infoLine.Value, nil
		}
	}

	return "", ErrNotFound
}

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
