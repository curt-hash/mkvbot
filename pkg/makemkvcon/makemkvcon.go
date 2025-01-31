package makemkvcon

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"iter"
	"log/slog"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type Config struct {
	ExePath          string ``                 // path to makemkvcon executable
	ProfilePath      string ``                 // path to makemkv profile XML
	ReadCacheSizeMB  int64  `validate:"min=1"` // --cache argument
	MinLengthSeconds int64  `validate:"min=1"` // --minlength argument
}

func (cfg *Config) Validate() error {
	if !fileExists(cfg.ExePath) {
		return fmt.Errorf("file %q not found", cfg.ExePath)
	}

	if cfg.ProfilePath != "" && !fileExists(cfg.ProfilePath) {
		return fmt.Errorf("file %q not found", cfg.ProfilePath)
	}

	return validate.Struct(cfg)
}

type MakeMKVCon struct {
	cfg *Config
}

func New(cfg *Config) (*MakeMKVCon, error) {
	if cfg.ExePath == "" {
		var err error
		if cfg.ExePath, err = FindExe(); err != nil {
			return nil, fmt.Errorf("find makemkvcon executable")
		}
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate %#+v: %w", cfg, err)
	}

	return &MakeMKVCon{
		cfg: cfg,
	}, nil
}

func (c *MakeMKVCon) ListDrives(ctx context.Context) (*LineIterator[[]*DriveScanLine], error) {
	// disc:9999 should trigger early termination since it is unlikely to exist.
	seq, err := c.runCmd(ctx, "info", "disc:9999")
	if err != nil {
		return nil, err
	}

	iter := &LineIterator[[]*DriveScanLine]{}
	iter.Seq = func(yield func(Line, error) bool) {
		for line, err := range seq {
			if !yield(line, err) {
				return
			}

			if err != nil {
				continue
			}

			if line.Kind() == LineKindDriveScan {
				l := line.(*DriveScanLine)
				if l.DriveName != "" {
					iter.result = append(iter.result, l)
				}
			}
		}
	}

	return iter, nil
}

func (c *MakeMKVCon) ScanDrive(ctx context.Context, driveIndex int) (*LineIterator[*Disc], error) {
	seq, err := c.runCmd(ctx, "info", fmt.Sprintf("disc:%d", driveIndex))
	if err != nil {
		return nil, err
	}

	d := &Disc{}
	iter := &LineIterator[*Disc]{
		result: d,
	}

	iter.Seq = func(yield func(Line, error) bool) {
		for line, err := range seq {
			if !yield(line, err) {
				return
			}

			if err != nil {
				continue
			}

			switch l := line.(type) {
			case *DiscInfoLine:
				d.Info = append(d.Info, l.InfoLine)
			case *TitleInfoLine:
				for l.TitleIndex >= len(d.Titles) {
					d.Titles = append(d.Titles, &Title{
						Index: l.TitleIndex,
					})
				}

				t := d.Titles[l.TitleIndex]
				t.Info = append(t.Info, l.InfoLine)
			case *StreamInfoLine:
				for l.TitleIndex >= len(d.Titles) {
					d.Titles = append(d.Titles, &Title{
						Index: l.TitleIndex,
					})
				}

				t := d.Titles[l.TitleIndex]
				for l.StreamIndex >= len(t.Streams) {
					t.Streams = append(t.Streams, &Stream{
						Index: l.StreamIndex,
					})
				}

				s := t.Streams[l.StreamIndex]
				s.Info = append(s.Info, l.InfoLine)
			}
		}
	}

	return iter, nil
}

func (c *MakeMKVCon) BackupTitle(ctx context.Context, driveIndex, titleIndex int, dstDir string) (iter.Seq2[Line, error], error) {
	if err := os.MkdirAll(dstDir, 0775); err != nil {
		return nil, fmt.Errorf("make directory %q: %w", dstDir, err)
	}

	return c.runCmd(
		ctx,
		"mkv",
		"--decrypt",
		fmt.Sprintf("--cache=%d", c.cfg.ReadCacheSizeMB),
		"--noscan",
		"--progress=-same",
		fmt.Sprintf("disc:%d", driveIndex),
		strconv.Itoa(titleIndex),
		dstDir,
	)
}

func (c *MakeMKVCon) runCmd(ctx context.Context, args ...string) (iter.Seq2[Line, error], error) {
	defaultArgs := []string{
		fmt.Sprintf("--minlength=%d", c.cfg.MinLengthSeconds),
		"-r",
	}

	if c.cfg.ProfilePath != "" {
		defaultArgs = append(defaultArgs, fmt.Sprintf("--profile=%s", c.cfg.ProfilePath))
	}

	cmd := exec.CommandContext(ctx, c.cfg.ExePath, slices.Concat(defaultArgs, args)...)
	cmd.WaitDelay = time.Second

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	slog.Debug("running command", "cmd", cmd.String())
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return func(yield func(Line, error) bool) {
		for line, err := range ParseLines(stdout) {
			if !yield(line, err) {
				return
			}
		}

		if err := cmd.Wait(); err != nil {
			yield(nil, err)
		}
	}, nil
}

func ParseLines(r io.Reader) iter.Seq2[Line, error] {
	return func(yield func(Line, error) bool) {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			s := scanner.Text()
			line, err := ParseLine(s)
			if err != nil {
				err = fmt.Errorf("parse line %q: %w", s, err)
			}
			if !yield(line, err) {
				return
			}
		}

		if err := scanner.Err(); err != nil {
			yield(nil, fmt.Errorf("scan lines from stdout: %w", err))
		}
	}
}

type LineIterator[T any] struct {
	Seq    iter.Seq2[Line, error]
	result T
	err    error
}

func (li *LineIterator[T]) GetResult() (T, error) {
	return li.result, li.err
}
