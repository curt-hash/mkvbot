package makemkvcon

type CurrentProgressLine struct {
	*progressLine
}

var _ Line = (*CurrentProgressLine)(nil)

func (l *CurrentProgressLine) Kind() LineKind {
	return LineKindCurrentProgress
}

func ParseCurrentProgressLine(s string) (*CurrentProgressLine, error) {
	l, err := parseProgressLine(s)
	if err != nil {
		return nil, err
	}

	return &CurrentProgressLine{l}, nil
}
