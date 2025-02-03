package makemkv

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

// Config is the makemkvcon configuration.
type Config struct {
	// ExePath is the path to the makemkvcon executable. It must exist.
	ExePath string

	// ProfilePath is the path to a makemkv profile XML file. makemkvcon relies
	// on it for the app_DefaultSelectionString setting, which determines what
	// streams (video, audio, and subtitles) are selected by default. It must
	// exist if non-empty.
	ProfilePath string

	// ReadCacheSizeMB is the value that is passed with the --cache argument to
	// makemkvcon. It must be at least 1.
	ReadCacheSizeMB int64 `validate:"min=1"`

	// MinLengthSeconds is the value that is passed with the --minlength argument
	// to makemkvcon. It must be at least 1.
	//
	// It filters out titles with video streams less than the given length, which
	// is very useful for weeding out unimportant streams.
	MinLengthSeconds int64 `validate:"min=1"`
}

// Validate returns an error if the configuration is invalid.
func (cfg *Config) Validate() error {
	if !fileExists(cfg.ExePath) {
		return fmt.Errorf("file %q not found", cfg.ExePath)
	}

	if cfg.ProfilePath != "" && !fileExists(cfg.ProfilePath) {
		return fmt.Errorf("file %q not found", cfg.ProfilePath)
	}

	return validate.Struct(cfg)
}

// Con is the interface for running makemkvcon commands.
type Con struct {
	cfg *Config

	defaultArgs []string
}

// New returns a new Con.
//
// If cfg.ExePath is empty, it will attempt to locate the executable
// automatically.
func New(cfg *Config) (*Con, error) {
	if cfg.ExePath == "" {
		var err error
		if cfg.ExePath, err = FindExe(); err != nil {
			return nil, fmt.Errorf("find makemkvcon executable")
		}
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate %#+v: %w", cfg, err)
	}

	defaultArgs := []string{
		fmt.Sprintf("--minlength=%d", cfg.MinLengthSeconds),
		"-r",
	}

	if cfg.ProfilePath != "" {
		defaultArgs = append(defaultArgs, fmt.Sprintf("--profile=%s", cfg.ProfilePath))
	}

	return &Con{
		cfg:         cfg,
		defaultArgs: defaultArgs,
	}, nil
}

// ListDrives returns the list of drives detected by makemkvcon.
func (c *Con) ListDrives(ctx context.Context) (*LineIterator[[]*DriveScanLine], error) {
	// disc:9999 should trigger early termination since it is unlikely to exist.
	seq, err := c.RunDefaultCmd(ctx, "info", "disc:9999")
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

// ScanDrive returns information about the disc in the given drive. The
// driveIndex should be obtained from ListDrives.
func (c *Con) ScanDrive(ctx context.Context, driveIndex int) (*LineIterator[*Disc], error) {
	seq, err := c.RunDefaultCmd(ctx, "info", fmt.Sprintf("disc:%d", driveIndex))
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
				t := d.GetTitle(l.TitleIndex)
				t.Info = append(t.Info, l.InfoLine)
			case *StreamInfoLine:
				s := d.GetTitle(l.TitleIndex).GetStream(l.StreamIndex)
				s.Info = append(s.Info, l.InfoLine)
			}
		}
	}

	return iter, nil
}

// BackupTitle creates a backup of title titleIndex of drive driveIndex in
// dstDir. The directory is created automatically if necessary.
func (c *Con) BackupTitle(ctx context.Context, driveIndex, titleIndex int, dstDir string) (iter.Seq2[Line, error], error) {
	if err := os.MkdirAll(dstDir, 0775); err != nil {
		return nil, fmt.Errorf("make directory %q: %w", dstDir, err)
	}

	return c.RunDefaultCmd(
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

// RunDefaultCmd calls RunCmd with default args in addition to the specified
// args. Default args include -r (machine-readable output), --minlength, and
// --profile.
func (c *Con) RunDefaultCmd(ctx context.Context, args ...string) (iter.Seq2[Line, error], error) {
	return c.RunCmd(ctx, slices.Concat(c.defaultArgs, args)...)
}

// RunCmd runs an arbitrary makemkvcon command with the given args. It
// terminates when the context is canceled or the command terminates.
func (c *Con) RunCmd(ctx context.Context, args ...string) (iter.Seq2[Line, error], error) {
	cmd := exec.CommandContext(ctx, c.cfg.ExePath, args...)
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

// ParseLines parses makemkvcon output lines from r. It returns a sequence of
// [Line, error] where either Line is a parsed line or err is non-nil. The
// sequence ends after all lines have been parsed and r returns EOF. Individual
// line parsing errors do not trigger an early return.
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

// LineIterator is a generic type that represents the lines output by a
// makemkvcon command and the generic final result.
type LineIterator[T any] struct {
	Seq    iter.Seq2[Line, error]
	result T
	err    error
}

// GetResult returns the final result of the command.
func (li *LineIterator[T]) GetResult() (T, error) {
	return li.result, li.err
}
