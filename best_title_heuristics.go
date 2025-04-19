package main

import (
	"log/slog"

	"github.com/curt-hash/mkvbot/pkg/makemkv"
)

type bestTitleHeuristic struct {
	name   string
	f      func(*makemkv.Disc) []*makemkv.Title
	weight int
}

var bestTitleHeuristics = []*bestTitleHeuristic{
	{
		name: "longest",
		f: func(d *makemkv.Disc) []*makemkv.Title {
			return d.TitlesWithLongestDuration()
		},
		weight: 1000,
	},
	{
		name: "most_chapters",
		f: func(d *makemkv.Disc) []*makemkv.Title {
			return d.TitlesWithMostChapters()
		},
		weight: 200,
	},
	{
		name: "angle_one",
		f: func(d *makemkv.Disc) []*makemkv.Title {
			return d.TitlesWithAngle(1)
		},
		weight: 300,
	},
	{
		name: "most_streams",
		f: func(d *makemkv.Disc) []*makemkv.Title {
			return d.TitlesWithMostStreams()
		},
		weight: 100,
	},
}

func findBestTitle(disc *makemkv.Disc, weights map[string]int) []*makemkv.Title {
	scores := make([]int, len(disc.Titles))
	for _, h := range bestTitleHeuristics {
		for _, title := range h.f(disc) {
			scores[title.Index] += weights[h.name]
		}
	}
	slog.Debug("scored titles", "scores", scores)

	return makemkv.Maximums(disc.Titles, func(title *makemkv.Title) (int, error) {
		return scores[title.Index], nil
	})
}
