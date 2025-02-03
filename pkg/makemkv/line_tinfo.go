package makemkv

// TitleInfoLine represents a makemkvcon "TINFO" output line, which provides
// information about a title.
//
// See https://makemkv.com/developers/usage.txt.
type TitleInfoLine struct {
	TitleIndex int

	*InfoLine
}

var _ Line = (*TitleInfoLine)(nil)

func (l *TitleInfoLine) Kind() LineKind {
	return LineKindTitleInfo
}

// ParseTitleInfoLine parses the string that follows "TINFO:" in the output of
// makemkvcon.
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
