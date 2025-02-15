package makemkv

import (
	"cmp"
	"time"

	"github.com/curt-hash/mkvbot/pkg/makemkv/defs"
)

// Disc is a sequence of titles plus some metadata.
type Disc struct {
	Info

	Titles []*Title
}

// GetTitle returns the title with the given index, creating it (and previous
// titles) if necessary.
func (d *Disc) GetTitle(index int) *Title {
	for index >= len(d.Titles) {
		d.Titles = append(d.Titles, &Title{
			Index: len(d.Titles),
		})
	}

	return d.Titles[index]
}

// TitleCount returns the number of titles on the disc.
func (d *Disc) TitleCount() int {
	return len(d.Titles)
}

// TitlesWithLongestDuration returns all titles that tie for maximum duration.
func (d *Disc) TitlesWithLongestDuration() []*Title {
	return Maximums(d.Titles, func(title *Title) (time.Duration, error) {
		return title.GetAttrDuration(defs.Duration)
	})
}

// TitlesWithAngle returns all titles with the given angle.
func (d *Disc) TitlesWithAngle(targetAngle int) []*Title {
	var matches []*Title
	for _, title := range d.Titles {
		if angle, err := title.GetAttrInt(defs.AngleInfo); err == nil && angle == targetAngle {
			matches = append(matches, title)
		}
	}

	return matches
}

// TitlesWithMostChapters returns all titles that tie for maximum number of
// chapters.
func (d *Disc) TitlesWithMostChapters() []*Title {
	return Maximums(d.Titles, func(title *Title) (int, error) {
		return title.GetAttrInt(defs.ChapterCount)
	})
}

// TitlesWithMostStreams returns all titles that tie for maximum number of
// streams.
func (d *Disc) TitlesWithMostStreams() []*Title {
	return Maximums(d.Titles, func(title *Title) (int, error) {
		return len(title.Streams), nil
	})
}

// Maximums returns all elements of the slice that maximize the given function,
// i.e., where f(e) = max(f(e0), f(e1), ..., f(eN)).
func Maximums[S []E, E any, V cmp.Ordered](s S, f func(E) (V, error)) S {
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
