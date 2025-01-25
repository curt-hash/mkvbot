package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"
)

func setDefaultLogger(w io.Writer, debug bool) {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	slog.SetDefault(slog.New(newLogHandler(w, level)))
}

type logHandler struct {
	w     io.Writer
	level slog.Leveler
}

var _ slog.Handler = (*logHandler)(nil)

func newLogHandler(w io.Writer, level slog.Leveler) *logHandler {
	return &logHandler{
		w:     w,
		level: level,
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

	_, err := h.w.Write(buf.Bytes())
	return err
}

func (h *logHandler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

func (h *logHandler) WithGroup(string) slog.Handler {
	return h
}
