package makemkvcon

type TotalProgressLine struct {
	*progressLine
}

var _ Line = (*TotalProgressLine)(nil)

func (l *TotalProgressLine) Kind() LineKind {
	return LineKindTotalProgress
}

func ParseTotalProgressLine(s string) (*TotalProgressLine, error) {
	l, err := parseProgressLine(s)
	if err != nil {
		return nil, err
	}

	return &TotalProgressLine{l}, nil
}
