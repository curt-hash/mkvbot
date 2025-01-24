package makemkvcon

import (
	"time"

	"github.com/curt-hash/mkvbot/pkg/makemkvcon/defs"
	"golang.org/x/exp/constraints"
)

type Disc struct {
	Info

	Titles []*Title
}

func (d *Disc) TitleCount() int {
	return len(d.Titles)
}

func (d *Disc) TitlesWithLongestDuration() []*Title {
	return Maximums(d.Titles, func(title *Title) (time.Duration, error) {
		return title.GetAttrDuration(defs.Duration)
	})
}

func (d *Disc) TitlesWithAngle(targetAngle int) []*Title {
	var matches []*Title
	for _, title := range d.Titles {
		if angle, err := title.GetAttrInt(defs.AngleInfo); err == nil && angle == targetAngle {
			matches = append(matches, title)
		}
	}

	return matches
}

func (d *Disc) TitlesWithMostChapters() []*Title {
	return Maximums(d.Titles, func(title *Title) (int, error) {
		return title.GetAttrInt(defs.ChapterCount)
	})
}

func (d *Disc) TitlesWithMostStreams() []*Title {
	return Maximums(d.Titles, func(title *Title) (int, error) {
		return len(title.Streams), nil
	})
}

func Maximums[S []E, E any, V constraints.Ordered](s S, f func(E) (V, error)) S {
	var (
		maxV     V
		maximums S
	)
	for _, e := range s {
		v, err := f(e)
		if err != nil {
			continue
		}

		switch {
		case v > maxV:
			maxV = v
			maximums = S{e}
		case v == maxV:
			maximums = append(maximums, e)
		}
	}

	return maximums
}
