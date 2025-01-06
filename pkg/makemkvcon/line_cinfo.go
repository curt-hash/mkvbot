package makemkvcon

type DiscInfoLine struct {
	*InfoLine
}

var _ Line = (*DiscInfoLine)(nil)

func (l *DiscInfoLine) Kind() LineKind {
	return LineKindDiscInfo
}

func ParseDiscInfoLine(s string) (*DiscInfoLine, error) {
	l, err := parseInfoLine(0, s)
	if err != nil {
		return nil, err
	}

	return &DiscInfoLine{l}, nil
}
