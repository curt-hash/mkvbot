package moviedb

// MovieMetadata is metadata about a movie.
//
// It contains the fields necessary for Plex's file naming scheme:
// https://support.plex.tv/articles/naming-and-organizing-your-movie-media-files
type MovieMetadata struct {
	Name string
	Year int

	// ID is the movie database identifier, e.g., "imdb-tt0118715".
	ID string
}

// MovieDB is the interface implemented by movie databases such as IMDb.
type MovieDB interface {
	// SearchMovies returns a list of results matching query q (typically the
	// movie title) ordered by relevance.
	SearchMovies(q string) ([]*MovieMetadata, error)
}
