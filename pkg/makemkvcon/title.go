package makemkvcon

type Title struct {
	Index int

	Info

	Streams []*Stream
}
