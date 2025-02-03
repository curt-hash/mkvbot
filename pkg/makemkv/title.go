package makemkv

// Title is a collection of Streams plus some metadata. It is identified by an
// index number. A Disc is made up of multiple Titles.
//
// Streams may be video, audio, or subtitles.
type Title struct {
	// Index is the index given by makemkv. Title numbers appear to be
	// deterministic if makemkv is run with the same --minlength argument.
	Index int

	Info

	Streams []*Stream
}

// GetStream returns the stream with the given index, creating it (and all
// prior streams) as necessary.
func (t *Title) GetStream(index int) *Stream {
	for index >= len(t.Streams) {
		t.Streams = append(t.Streams, &Stream{
			Index: len(t.Streams),
		})
	}

	return t.Streams[index]
}
