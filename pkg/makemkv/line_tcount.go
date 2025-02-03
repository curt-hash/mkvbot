package makemkv

import (
	"fmt"
	"strconv"
)

// TitleCountLine represents a "TCOUNT" makemkvcon output line, which describes
// the number of titles found on a disc.
//
// See https://makemkv.com/developers/usage.txt.
type TitleCountLine struct {
	Count int
}

var _ Line = (*TitleCountLine)(nil)

func (l *TitleCountLine) Kind() LineKind {
	return LineKindTitleCount
}

// ParseTitleCountLine parses the string that follows "TCOUNT:" in the output
// of makemkvcon.
func ParseTitleCountLine(s string) (*TitleCountLine, error) {
	count, err := strconv.Atoi(s)
	if err != nil {
		return nil, fmt.Errorf("expected %q to be an integer: %w", s, err)
	}

	return &TitleCountLine{
		Count: count,
	}, nil
}
