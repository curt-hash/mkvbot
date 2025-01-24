package makemkvcon

type TitleInfoLine struct {
	TitleIndex int `json:"title_index"`

	*InfoLine
}

var _ Line = (*TitleInfoLine)(nil)

func (l *TitleInfoLine) Kind() LineKind {
	return LineKindTitleInfo
}

func ParseTitleInfoLine(s string) (*TitleInfoLine, error) {
	l, err := parseInfoLine(1, s)
	if err != nil {
		return nil, err
	}

	return &TitleInfoLine{
		TitleIndex: l.prefix[0],
		InfoLine:   l,
	}, nil
}
