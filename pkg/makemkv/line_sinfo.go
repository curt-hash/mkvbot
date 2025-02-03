package makemkv

// StreamInfoLine represents an "SINFO" makemkvcon output line, which provides
// information about a stream.
//
// See https://makemkv.com/developers/usage.txt.
type StreamInfoLine struct {
	TitleIndex  int
	StreamIndex int

	*InfoLine
}

var _ Line = (*StreamInfoLine)(nil)

func (l *StreamInfoLine) Kind() LineKind {
	return LineKindStreamInfo
}

// ParseStreamInfoLine parses the string that follows "SINFO:" in the output of
// makemkvcon.
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
