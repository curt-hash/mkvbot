package makemkv

// TotalProgressLine represents a makemkvcon "PRGT" output line, which
// describes the overall task being performed.
//
// See https://makemkv.com/developers/usage.txt.
type TotalProgressLine struct {
	*progressLine
}

var _ Line = (*TotalProgressLine)(nil)

func (l *TotalProgressLine) Kind() LineKind {
	return LineKindTotalProgress
}

// ParseTotalProgressLine parses the string that follows "PRGT:" in the output
// of makemkvcon.
func ParseTotalProgressLine(s string) (*TotalProgressLine, error) {
	l, err := parseProgressLine(s)
	if err != nil {
		return nil, err
	}

	return &TotalProgressLine{l}, nil
}
