package makemkv

// DiscInfoLine represents a makemkvcon "CINFO" output line, which provides
// information about a disc.
//
// See https://makemkv.com/developers/usage.txt.
type DiscInfoLine struct {
	*InfoLine
}

var _ Line = (*DiscInfoLine)(nil)

func (l *DiscInfoLine) Kind() LineKind {
	return LineKindDiscInfo
}

// ParseDiscInfo parses the string that follows "CINFO:" in the output of
// makemkvcon.
func ParseDiscInfoLine(s string) (*DiscInfoLine, error) {
	l, err := parseInfoLine(0, s)
	if err != nil {
		return nil, err
	}

	return &DiscInfoLine{l}, nil
}
