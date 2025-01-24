package moviedb

import (
	"fmt"
	"net/http"

	"github.com/StalkR/imdb"
)

type IMDb struct {
	client *http.Client
}

var _ DB = (*IMDb)(nil)

func NewIMDb() *IMDb {
	return &IMDb{
		client: &http.Client{
			Transport: &customTransport{http.DefaultTransport},
		},
	}
}

func (s *IMDb) FuzzySearchTitle(q string) ([]*Metadata, error) {
	titles, err := imdb.SearchTitle(s.client, q)
	if err != nil {
		return nil, fmt.Errorf("search imdb: %w", err)
	}

	results := make([]*Metadata, len(titles))
	for i, title := range titles {
		results[i] = &Metadata{
			Name: title.Name,
			Year: title.Year,
			Tag:  fmt.Sprintf("imdb-%s", title.ID),
		}
	}

	return results, nil
}

type customTransport struct {
	http.RoundTripper
}

func (t *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36")
	return t.RoundTripper.RoundTrip(req)
}
