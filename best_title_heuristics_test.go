package main

import (
	"testing"

	"github.com/curt-hash/mkvbot/pkg/makemkv"
	"github.com/curt-hash/mkvbot/pkg/makemkv/defs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// makeDisc builds a Disc with titles described by the provided specs.
// Each spec is a map of defs.Attr to string value, plus an optional
// "__streams__" key (int) for stream count.
func makeDisc(specs []map[defs.Attr]string) *makemkv.Disc {
	d := &makemkv.Disc{}
	for i, spec := range specs {
		t := d.GetTitle(i)
		for attr, val := range spec {
			t.Info = append(t.Info, &makemkv.Attribute{
				ID:    int(attr),
				Value: makemkv.Str(val),
			})
		}
	}
	return d
}

// addStreams appends n empty streams to title at index i.
func addStreams(d *makemkv.Disc, i, n int) {
	t := d.GetTitle(i)
	for j := range n {
		t.GetStream(j)
	}
}

func defaultWeights() map[string]int64 {
	weights := make(map[string]int64)
	for _, h := range bestTitleHeuristics {
		weights[h.name] = h.weight
	}
	return weights
}

func TestFindBestTitle(t *testing.T) {
	t.Run("longest wins", func(t *testing.T) {
		d := makeDisc([]map[defs.Attr]string{
			{defs.Duration: "1:00:00"},
			{defs.Duration: "2:00:00"},
			{defs.Duration: "0:30:00"},
		})
		best := findBestTitle(d, defaultWeights())
		require.Len(t, best, 1)
		assert.Equal(t, 1, best[0].Index)
	})

	t.Run("tie resolved by secondary heuristic", func(t *testing.T) {
		// Two titles with the same duration; title 1 has more chapters.
		d := makeDisc([]map[defs.Attr]string{
			{defs.Duration: "2:00:00", defs.ChapterCount: "10"},
			{defs.Duration: "2:00:00", defs.ChapterCount: "20"},
		})
		best := findBestTitle(d, defaultWeights())
		require.Len(t, best, 1)
		assert.Equal(t, 1, best[0].Index)
	})

	t.Run("full tie returns all", func(t *testing.T) {
		d := makeDisc([]map[defs.Attr]string{
			{defs.Duration: "2:00:00", defs.ChapterCount: "10", defs.AngleInfo: "1"},
			{defs.Duration: "2:00:00", defs.ChapterCount: "10", defs.AngleInfo: "1"},
		})
		addStreams(d, 0, 5)
		addStreams(d, 1, 5)
		best := findBestTitle(d, defaultWeights())
		assert.Len(t, best, 2)
	})

	t.Run("zero weight heuristic has no effect", func(t *testing.T) {
		// Title 0 is longer; title 1 has more chapters.
		// Zero out the most_chapters weight so title 0 should still win.
		d := makeDisc([]map[defs.Attr]string{
			{defs.Duration: "2:00:00", defs.ChapterCount: "5"},
			{defs.Duration: "1:00:00", defs.ChapterCount: "20"},
		})
		weights := defaultWeights()
		weights["most_chapters"] = 0
		best := findBestTitle(d, weights)
		require.Len(t, best, 1)
		assert.Equal(t, 0, best[0].Index)
	})

	t.Run("single title always wins", func(t *testing.T) {
		d := makeDisc([]map[defs.Attr]string{
			{defs.Duration: "1:30:00"},
		})
		best := findBestTitle(d, defaultWeights())
		require.Len(t, best, 1)
		assert.Equal(t, 0, best[0].Index)
	})
}
