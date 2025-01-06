package makemkvcon

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

// From https://makemkv.com/developers/usage.txt:
//
// MSG:code,flags,count,message,format,param0,param1,...
// code - unique message code, should be used to identify particular string in language-neutral way.
// flags - message flags, see AP_UIMSG_xxx flags in apdefs.h
// count - number of parameters
// message - raw message string suitable for output
// format - format string used for message. This string is localized and subject to change, unlike message code.
// paramX - parameter for message
type MessageLine struct {
	Code      int      `json:"code"`
	Flags     int      `json:"flags"`
	NumParams int      `json:"num_params"`
	Message   string   `json:"message"`
	Format    string   `json:"format"`
	Params    []string `json:"params"`
}

var _ Line = (*MessageLine)(nil)

func (l *MessageLine) Kind() LineKind {
	return LineKindMessage
}

func ParseMessageLine(s string) (*MessageLine, error) {
	tokens := strings.SplitN(s, ",", 4)
	if len(tokens) != 4 {
		return nil, fmt.Errorf("expected 4 tokens, got %d", len(tokens))
	}

	numParams, err := strconv.Atoi(tokens[2])
	if err != nil {
		return nil, fmt.Errorf("parse count %q: %w", tokens[2], err)
	}

	reader := csv.NewReader(strings.NewReader(s))
	reader.FieldsPerRecord = 5 + numParams

	record, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("parse csv line %q: %w", s, err)
	}

	l := &MessageLine{
		NumParams: numParams,
		Message:   record[3],
		Format:    record[4],
		Params:    record[5:],
	}

	if l.Code, err = strconv.Atoi(record[0]); err != nil {
		return nil, fmt.Errorf("parse code %q: %w", record[0], err)
	}

	if l.Flags, err = strconv.Atoi(record[1]); err != nil {
		return nil, fmt.Errorf("parse flags %q: %w", record[1], err)
	}

	return l, nil
}
