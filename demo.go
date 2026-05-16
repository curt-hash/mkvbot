package main

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"time"

	"github.com/curt-hash/mkvbot/pkg/makemkv"
	"github.com/curt-hash/mkvbot/pkg/makemkv/defs"
	"github.com/curt-hash/mkvbot/pkg/moviedb"
	"golang.org/x/sync/errgroup"
)

func runDemo(ctx context.Context) error {
	tui := newTextUserInterface(newBeeper(false))
	setDefaultLogger([]io.Writer{tui.logBox}, false)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		tui.waitForInterrupt()
		cancel()
	}()

	var tasks errgroup.Group
	tasks.Go(tui.run)

	go func() {
		time.Sleep(50 * time.Millisecond)
		populateDemoData(tui)
	}()

	<-ctx.Done()
	tui.Stop()
	return errors.Join(nil, tasks.Wait())
}

func attr(id defs.Attr, value string) *makemkv.Attribute {
	return &makemkv.Attribute{ID: int(id), Value: makemkv.Str(value)}
}

func attrCode(code defs.TypeCode) *makemkv.Attribute {
	return &makemkv.Attribute{ID: int(defs.Type), Code: int(code)}
}

func populateDemoData(tui *textUserInterface) {
	tui.setDriveInfo("BD-RE ACME OPT-9000 2.01 000000000000", "D:")

	tui.setDiscInfo(makemkv.Info{
		attr(defs.Type, "Blu-ray disc"),
		attr(defs.Name, "SAMPLE_DISC"),
		attr(defs.TreeInfo, "SAMPLE_DISC"),
		attr(defs.VolumeName, "SAMPLE_DISC"),
		attr(defs.OrderWeight, "0"),
	})

	tui.setMovieMetadata(&moviedb.MovieMetadata{
		Name: "Movie",
		Year: 2021,
		ID:   "imdb-tt0000000",
	})

	title := &makemkv.Title{Index: 1}
	title.Info = makemkv.Info{
		attr(defs.ChapterCount, "24"),
		attr(defs.Duration, "2:07:00"),
		attr(defs.DiskSize, "38,654,705,664"),
		attr(defs.SourceFileName, "00800.mpls"),
		attr(defs.SegmentsMap, "1:1"),
		attr(defs.OutputFileName, "00800.mkv"),
	}

	video := title.GetStream(0)
	video.Info = makemkv.Info{
		attrCode(defs.TypeCodeVideo),
		attr(defs.TreeInfo, "H.265/HEVC Video"),
		attr(defs.VideoSize, "1920x1080"),
		attr(defs.Bitrate, "35000 Kb/s"),
	}

	audio1 := title.GetStream(1)
	audio1.Info = makemkv.Info{
		attrCode(defs.TypeCodeAudio),
		attr(defs.CodecLong, "DTS-HD Master Audio"),
		attr(defs.AudioChannelLayoutName, "7.1"),
		attr(defs.LangName, "English"),
	}

	audio2 := title.GetStream(2)
	audio2.Info = makemkv.Info{
		attrCode(defs.TypeCodeAudio),
		attr(defs.CodecLong, "DTS"),
		attr(defs.AudioChannelLayoutName, "5.1(side)"),
		attr(defs.LangName, "English"),
	}

	sub1 := title.GetStream(3)
	sub1.Info = makemkv.Info{
		attrCode(defs.TypeCodeSubtitles),
		attr(defs.LangName, "English"),
	}

	sub2 := title.GetStream(4)
	sub2.Info = makemkv.Info{
		attrCode(defs.TypeCodeSubtitles),
		attr(defs.LangName, "Spanish"),
	}

	tui.setTitleInfo(title)

	tui.setStatus("Backing up title to /movies/Movie (2021) {imdb-tt0000000}")
	tui.setTask("Saving all titles to MKV files")
	tui.setSubtask("Saving to MKV file")
	tui.setProgress(0.42)

	for _, m := range []string{
		"MakeMKV v1.18.0 started",
		"Loaded content-hash table, will verify integrity of HDTS files.",
		"File 00001.mpls was added as title #1",
		"File 00800.mpls was added as title #2",
		"Title 00001.mpls has length of 92 seconds which is less than minimum title length of 1800 seconds and was skipped",
		"Title 00002.mpls has length of 45 seconds which is less than minimum title length of 1800 seconds and was skipped",
		"Title 00003.mpls has length of 120 seconds which is less than minimum title length of 1800 seconds and was skipped",
		"Saving 1 titles into directory File:///movies/Movie (2021) {imdb-tt0000000}",
	} {
		slog.Info(m, "source", "makemkv")
	}
}
