package main

import (
	"log/slog"

	"github.com/curt-hash/mkvbot/pkg/makemkvcon"
)

type bestTitleHeuristic struct {
	name      string
	f         func(*makemkvcon.Disc) []*makemkvcon.Title
	weight    int64
	flagName  string
	flagUsage string
}

var bestTitleHeuristics = []*bestTitleHeuristic{
	{
		name: "longest",
		f: func(d *makemkvcon.Disc) []*makemkvcon.Title {
			return d.TitlesWithLongestDuration()
		},
		weight:    1000,
		flagName:  "longest-title-weight",
		flagUsage: "`WEIGHT` given to longest title(s)",
	},
	{
		name: "most chapters",
		f: func(d *makemkvcon.Disc) []*makemkvcon.Title {
			return d.TitlesWithMostChapters()
		},
		weight:    200,
		flagName:  "most-chapters-weight",
		flagUsage: "`WEIGHT` given to title(s) with the most chapters",
	},
	{
		name: "angle one",
		f: func(d *makemkvcon.Disc) []*makemkvcon.Title {
			return d.TitlesWithAngle(1)
		},
		weight:    300,
		flagName:  "angle-one-weight",
		flagUsage: "`WEIGHT` given to title(s) with angle one",
	},
	{
		name: "most streams",
		f: func(d *makemkvcon.Disc) []*makemkvcon.Title {
			return d.TitlesWithMostStreams()
		},
		weight:    100,
		flagName:  "most-streams-weight",
		flagUsage: "`WEIGHT` given to title(s) with the most streams",
	},
}

func findBestTitle(disc *makemkvcon.Disc, weights map[string]int64) []*makemkvcon.Title {
	scores := make([]int64, len(disc.Titles))
	for _, h := range bestTitleHeuristics {
		for _, title := range h.f(disc) {
			scores[title.Index] += weights[h.name]
		}
	}
	slog.Debug("scored titles", "scores", scores)

	return makemkvcon.Maximums(disc.Titles, func(title *makemkvcon.Title) (int64, error) {
		return scores[title.Index], nil
	})
}
