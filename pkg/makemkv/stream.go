package makemkv

import (
	"github.com/curt-hash/mkvbot/pkg/makemkv/defs"
)

// Stream is a video, audio or subtitles stream. A Title is made up of multiple
// Streams.
type Stream struct {
	// Index is the index given by makemkv.
	Index int

	Info
}

// Type returns the type code of the stream.
func (s *Stream) Type() defs.TypeCode {
	return defs.TypeCode(s.GetCodeDefault(defs.Type, 0))
}
