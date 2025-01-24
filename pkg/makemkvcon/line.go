package makemkvcon

import (
	"fmt"
	"strings"
)

type LineKind int

const (
	LineKindUnknown LineKind = iota
	LineKindMessage
	LineKindDriveScan
	LineKindTitleCount
	LineKindDiscInfo
	LineKindTitleInfo
	LineKindStreamInfo
	LineKindCurrentProgress
	LineKindTotalProgress
	LineKindProgressBar
)

type Line interface {
	Kind() LineKind
}

func ParseLine(s string) (Line, error) {
	before, after, found := strings.Cut(s, ":")
	if !found {
		return nil, fmt.Errorf("no colon found in line %q: %w", s, ErrUnhandledLine)
	}

	var (
		line Line
		err  error
	)
	switch before {
	case "MSG":
		line, err = ParseMessageLine(after)
	case "DRV":
		line, err = ParseDriveScanLine(after)
	case "TCOUNT":
		line, err = ParseTitleCountLine(after)
	case "CINFO":
		line, err = ParseDiscInfoLine(after)
	case "TINFO":
		line, err = ParseTitleInfoLine(after)
	case "SINFO":
		line, err = ParseStreamInfoLine(after)
	case "PRGC":
		line, err = ParseCurrentProgressLine(after)
	case "PRGT":
		line, err = ParseTotalProgressLine(after)
	case "PRGV":
		line, err = ParseProgressBarLine(after)
	default:
		return nil, fmt.Errorf("unhandled line prefix %q: %w", before, ErrUnhandledLine)
	}

	return line, err
}
