package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"
)

func setDefaultLogger(writers []io.Writer, debug bool) {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	slog.SetDefault(slog.New(newLogHandler(writers, level)))
}

type logHandler struct {
	writers []io.Writer
	level   slog.Leveler
}

var _ slog.Handler = (*logHandler)(nil)

func newLogHandler(writers []io.Writer, level slog.Leveler) *logHandler {
	return &logHandler{
		writers: writers,
		level:   level,
	}
}

func (h *logHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *logHandler) Handle(_ context.Context, r slog.Record) error {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s %-7s %s", r.Time.Format(time.TimeOnly), fmt.Sprintf("[%s]", r.Level), r.Message)
	numAttrs := r.NumAttrs()
	if numAttrs > 0 {
		buf.WriteString(" (")
	}
	i := 0
	r.Attrs(func(attr slog.Attr) bool {
		if i > 0 {
			buf.WriteByte(' ')
		}
		fmt.Fprintf(&buf, "%s=%v", attr.Key, attr.Value)
		i++
		return true
	})
	if numAttrs > 0 {
		buf.WriteByte(')')
	}
	buf.WriteByte('\n')

	var errs error
	for _, w := range h.writers {
		if _, err := w.Write(buf.Bytes()); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return errs
}

func (h *logHandler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

func (h *logHandler) WithGroup(string) slog.Handler {
	return h
}
