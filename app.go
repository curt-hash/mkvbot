package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/curt-hash/mkvbot/pkg/eject"
	"github.com/curt-hash/mkvbot/pkg/makemkv"
	"github.com/curt-hash/mkvbot/pkg/makemkv/defs"
	"github.com/curt-hash/mkvbot/pkg/moviedb"
	"github.com/gen2brain/beeep"
	"golang.org/x/sync/errgroup"
)

type (
	applicationConfig struct {
		outputDirPath              string
		makemkvConfig              *makemkv.Config
		debug                      bool
		quiet                      bool
		bestTitleHeuristicsWeights map[string]int64
		askForTitle                bool
		logFilePath                string
	}

	application struct {
		cfg     *applicationConfig
		con     *makemkv.Con
		tui     *textUserInterface
		logFile *os.File
	}
)

func (cfg *applicationConfig) validate() error {
	if _, err := os.Stat(cfg.outputDirPath); err != nil {
		return fmt.Errorf("stat %q: %w", cfg.outputDirPath, err)
	}

	for _, h := range bestTitleHeuristics {
		if _, ok := cfg.bestTitleHeuristicsWeights[h.name]; !ok {
			return fmt.Errorf("missing weight for best title heuristic: %q", h.name)
		}
	}

	return nil
}

func newApplication(cfg *applicationConfig) (*application, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validate config %#+v: %w", cfg, err)
	}

	con, err := makemkv.New(cfg.makemkvConfig)
	if err != nil {
		return nil, fmt.Errorf("initialize makemkv controller: %w", err)
	}

	tui := newTextUserInterface(newBeeper(!cfg.quiet))

	logWriters := []io.Writer{tui.logBox}
	var logFile *os.File
	if cfg.logFilePath != "" {
		if logFile, err = os.OpenFile(cfg.logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
			return nil, fmt.Errorf("open %q: %w", cfg.logFilePath, err)
		}

		logWriters = append(logWriters, logFile)
	}
	setDefaultLogger(logWriters, cfg.debug)

	return &application{
		cfg:     cfg,
		con:     con,
		tui:     tui,
		logFile: logFile,
	}, nil
}

func (app *application) run(ctx context.Context) (err error) {
	defer func() {
		if app.logFile != nil {
			err = errors.Join(err, app.logFile.Close())
		}
	}()

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		app.tui.waitForInterrupt()
		cancel()
	}()

	var tasks errgroup.Group
	tasks.Go(app.tui.run)

	err = app.doBackupLoop(ctx)
	app.tui.Stop()
	return errors.Join(err, tasks.Wait())
}

func (app *application) doBackupLoop(ctx context.Context) error {
	app.tui.setStatus("Inspecting drives")
	drive, err := app.getDrive(ctx)
	if err != nil {
		return err
	}

	app.tui.setDriveInfo(drive.DriveName.String(), drive.VolumeName.String())

	for ctx.Err() == nil {
		if err := app.tryBackupBestTitle(ctx, drive); err != nil {
			slog.Error(err.Error())
		}

		app.tui.setStatus("Sleeping for a moment")
		select {
		case <-ctx.Done():
		case <-time.After(2 * time.Second):
		}
	}

	return nil
}

func (app *application) getDrive(ctx context.Context) (*makemkv.DriveScan, error) {
	iter, err := app.con.ListDrives(ctx)
	if err != nil {
		return nil, fmt.Errorf("list drives: %w", err)
	}

	for line, err := range iter.Seq {
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		switch {
		case line.Message != nil:
			slog.Info(line.Message.Message.String(), "source", "makemkv")
		}
	}

	drives, err := iter.GetResult()
	if err != nil {
		return nil, fmt.Errorf("list drives: %w", err)
	}

	switch len(drives) {
	case 0:
		return nil, fmt.Errorf("no drives found")
	case 1:
	default:
		slog.Warn("multiple drives not supported yet; using first drive only")
	}

	return drives[0], nil
}

func (app *application) tryBackupBestTitle(ctx context.Context, drive *makemkv.DriveScan) error {
	defer func() {
		app.tui.setDiscInfo(nil)
		app.tui.setMovieMetadata(nil)
		app.tui.setTitleInfo(nil)
	}()

	app.tui.setStatus("Scanning drive %q", drive.VolumeName)
	iter, err := app.con.ScanDrive(ctx, drive.Index)
	if err != nil {
		return fmt.Errorf("scan drive %q: %w", drive.VolumeName, err)
	}

	for line, err := range iter.Seq {
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		switch {
		case line.CurrentTask != nil:
			app.tui.setTask("%s", line.CurrentTask.Task.Name)
		case line.CurrentSubtask != nil:
			app.tui.setSubtask("%s", line.CurrentSubtask.Task.Name)
		case line.Progress != nil:
			app.tui.setProgress(line.Progress.TaskProgress())
		case line.Message != nil:
			slog.Debug(line.Message.Message.String(), "source", "makemkv")
		}
	}

	disc, err := iter.GetResult()
	if err != nil {
		return err
	}

	if disc.TitleCount() == 0 {
		slog.Debug("no titles found")
		return nil
	}

	app.tui.setDiscInfo(disc.Info)

	app.tui.setStatus("Getting movie metadata")
	movieMetadata, err := app.getMovieMetadata(ctx, disc)
	if err != nil {
		return fmt.Errorf("get movie metadata: %w", err)
	}
	app.tui.setMovieMetadata(movieMetadata)
	fileName := makeFileName(movieMetadata)

	var (
		title *makemkv.Title
		best  []*makemkv.Title
	)
	app.tui.setStatus("Finding best title")
	if app.cfg.askForTitle {
		best = disc.Titles
	} else {
		best = findBestTitle(disc, app.cfg.bestTitleHeuristicsWeights)
	}
	switch len(best) {
	case 0:
		return fmt.Errorf("no best titles")
	case 1:
		title = best[0]
	default:
		if title, err = app.tui.getBestTitle(ctx, best); err != nil {
			return fmt.Errorf("get best title: %w", err)
		}
	}
	app.tui.setTitleInfo(title)

	app.tui.setStatus("Backing up title")
	if err := app.backupTitle(ctx, drive, title, fileName); err != nil {
		return fmt.Errorf("backup longest title: %w", err)
	}

	app.tui.setStatus("Ejecting disc")
	if err := eject.Eject(ctx, drive.VolumeName.String()); err != nil {
		return fmt.Errorf("eject disc: %w", err)
	}

	app.tui.beep()
	return nil
}

func (app *application) getMovieMetadata(ctx context.Context, disc *makemkv.Disc) (*moviedb.MovieMetadata, error) {
	name, err := disc.GetAttr(defs.Name)
	if err != nil {
		return nil, fmt.Errorf("get disc attr %s: %w", defs.Name, err)
	}

	q := regexp.MustCompile("[^a-zA-Z0-9 ]+").ReplaceAllString(name, " ")
	q, err = app.tui.getMovieTitleForSearch(ctx, q)
	if err != nil {
		return nil, err
	}

	metadata, err := searchMovieDB(q)
	if err != nil {
		slog.Warn("movie metadata lookup failed", "err", err)
		metadata = &moviedb.MovieMetadata{}
	}

	return app.tui.getMovieMetadata(ctx, metadata)
}

func (app *application) backupTitle(ctx context.Context, drive *makemkv.DriveScan, title *makemkv.Title, fileName string) error {
	dstDir := filepath.Join(app.cfg.outputDirPath, fileName)
	dstPath := filepath.Join(dstDir, fmt.Sprintf("%s.mkv", fileName))
	if _, err := os.Stat(dstPath); err == nil {
		return fmt.Errorf("output file exists: %q", dstPath)
	}

	app.tui.setStatus("Backing up title to %s", dstDir)
	seq, err := app.con.BackupTitle(ctx, drive.Index, title.Index, dstDir)
	if err != nil {
		return fmt.Errorf("backup title %d to %q: %w", title.Index, dstDir, err)
	}

	for line, err := range seq {
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		switch {
		case line.CurrentTask != nil:
			app.tui.setTask("%s", line.CurrentTask.Task.Name)
		case line.CurrentSubtask != nil:
			app.tui.setSubtask("%s", line.CurrentSubtask.Task.Name)
		case line.Progress != nil:
			app.tui.setProgress(line.Progress.TaskProgress())
		case line.Message != nil:
			slog.Info(line.Message.Message.String(), "source", "makemkv")
		}
	}

	name, err := title.GetAttr(defs.OutputFileName)
	if err != nil {
		return fmt.Errorf("title has no output file name")
	}

	expectedPath := filepath.Join(dstDir, name)
	if _, err := os.Stat(expectedPath); err != nil {
		return fmt.Errorf("backup file not found at expected path %q: %w", expectedPath, err)
	}

	if err := os.Rename(expectedPath, dstPath); err != nil {
		return fmt.Errorf("rename %q to %q: %w", expectedPath, dstPath, err)
	}

	return nil
}

func makeFileName(metadata *moviedb.MovieMetadata) string {
	return fmt.Sprintf("%s (%d) {%s}", sanitizeFileName(metadata.Name), metadata.Year, sanitizeFileName(metadata.ID))
}

func searchMovieDB(q string) (*moviedb.MovieMetadata, error) {
	results, err := moviedb.NewIMDb().SearchMovies(q)
	if err != nil {
		return nil, fmt.Errorf("search for %q: %w", q, err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results")
	}

	return results[0], nil
}

type beeper struct {
	enabled bool
}

func newBeeper(enabled bool) *beeper {
	return &beeper{
		enabled: enabled,
	}
}

func (b *beeper) beep() {
	if b.enabled {
		for range 2 {
			if err := beeep.Beep(2400, 500); err != nil {
				slog.Error("beep error", "err", err)
			}
		}
	}
}
