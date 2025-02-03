package makemkv

// CurrentProgressLine represents a makemkvcon "PRGC" output line, which
// describes the current sub-task.
//
// See https://makemkv.com/developers/usage.txt.
type CurrentProgressLine struct {
	*progressLine
}

var _ Line = (*CurrentProgressLine)(nil)

func (l *CurrentProgressLine) Kind() LineKind {
	return LineKindCurrentProgress
}

// ParseCurrentProgressLine parses the string that follows "PRGC:" in the
// output of makemkvcon.
func ParseCurrentProgressLine(s string) (*CurrentProgressLine, error) {
	l, err := parseProgressLine(s)
	if err != nil {
		return nil, err
	}

	return &CurrentProgressLine{l}, nil
}
