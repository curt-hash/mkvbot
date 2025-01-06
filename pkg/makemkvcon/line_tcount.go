package makemkvcon

import (
	"fmt"
	"strconv"
)

// From https://makemkv.com/developers/usage.txt:
//
// TCOUT:count
// count - titles count
type TitleCountLine struct {
	Count int `json:"count"`
}

var _ Line = (*TitleCountLine)(nil)

func (l *TitleCountLine) Kind() LineKind {
	return LineKindTitleCount
}

func ParseTitleCountLine(s string) (*TitleCountLine, error) {
	count, err := strconv.Atoi(s)
	if err != nil {
		return nil, fmt.Errorf("expected %q to be an integer: %w", s, err)
	}

	return &TitleCountLine{
		Count: count,
	}, nil
}
