package makemkvcon

type StreamInfoLine struct {
	TitleIndex  int `json:"title_index"`
	StreamIndex int `json:"stream_index"`

	*InfoLine
}

var _ Line = (*StreamInfoLine)(nil)

func (l *StreamInfoLine) Kind() LineKind {
	return LineKindStreamInfo
}

func ParseStreamInfoLine(s string) (*StreamInfoLine, error) {
	l, err := parseInfoLine(2, s)
	if err != nil {
		return nil, err
	}

	return &StreamInfoLine{
		TitleIndex:  l.prefix[0],
		StreamIndex: l.prefix[1],
		InfoLine:    l,
	}, nil
}
