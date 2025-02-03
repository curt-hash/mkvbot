package moviedb

import (
	"fmt"
	"net/http"

	"github.com/StalkR/imdb"
)

// IMDb interfaces with the imdb.com movie database.
type IMDb struct {
	client *http.Client
}

var _ MovieDB = (*IMDb)(nil)

// NewIMDb returns a new IMDb.
func NewIMDb() *IMDb {
	return &IMDb{
		client: &http.Client{
			Transport: &customTransport{http.DefaultTransport},
		},
	}
}

// SearchMovies implements MovieDB.
func (s *IMDb) SearchMovies(q string) ([]*MovieMetadata, error) {
	titles, err := imdb.SearchTitle(s.client, q)
	if err != nil {
		return nil, fmt.Errorf("search imdb: %w", err)
	}

	results := make([]*MovieMetadata, len(titles))
	for i, title := range titles {
		results[i] = &MovieMetadata{
			Name: title.Name,
			Year: title.Year,
			ID:   fmt.Sprintf("imdb-%s", title.ID),
		}
	}

	return results, nil
}

type customTransport struct {
	http.RoundTripper
}

func (t *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// The API used by StalkR/imdb denies requests based on headers, so mimic a
	// real browser. This is obviously somewhat fragile.
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36")
	return t.RoundTripper.RoundTrip(req)
}
