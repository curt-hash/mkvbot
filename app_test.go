package main

import (
	"regexp"
	"testing"

	"github.com/curt-hash/mkvbot/pkg/moviedb"
	"github.com/stretchr/testify/assert"
)

func TestMakeFileName(t *testing.T) {
	tests := []struct {
		name     string
		metadata *moviedb.MovieMetadata
		want     string
		wantRe   *regexp.Regexp
	}{
		{
			name:     "full metadata",
			metadata: &moviedb.MovieMetadata{Name: "Test Movie", Year: 2000, ID: "imdb-tt1234567"},
			want:     "Test Movie (2000) {imdb-tt1234567}",
		},
		{
			name:     "no ID",
			metadata: &moviedb.MovieMetadata{Name: "Another Film", Year: 2010},
			want:     "Another Film (2010)",
		},
		{
			name:     "year zero substituted with current year",
			metadata: &moviedb.MovieMetadata{Name: "Generic Title", Year: 0, ID: "imdb-tt0000000"},
			wantRe:   regexp.MustCompile(`^Generic Title \(\d{4}\) \{imdb-tt0000000\}$`),
		},
		{
			name:     "special characters in name",
			metadata: &moviedb.MovieMetadata{Name: "Test: Movie <Version>", Year: 2015, ID: "imdb-tt5555555"},
			want:     "Test- Movie Version (2015) {imdb-tt5555555}",
		},
		{
			name:     "no name falls back to timestamped unknown",
			metadata: &moviedb.MovieMetadata{Name: "", Year: 0, ID: ""},
			wantRe:   regexp.MustCompile(`^Unknown Title \d{4}-\d{2}-\d{2}T\d{2}-\d{2}-\d{2}Z \(\d{4}\)$`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeFileName(tt.metadata)
			if tt.wantRe != nil {
				assert.Regexp(t, tt.wantRe, got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
