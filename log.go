package examplegoslog

import (
	"context"
	"log/slog"
	"runtime"
	"time"
)

var _ slog.Handler = (*handler)(nil)

func New(h slog.Handler, fs ...f) *handler {
	return &handler{
		h:  h,
		fs: fs,
	}
}

type f func(ctx context.Context) slog.Attr

type handler struct {
	h  slog.Handler
	fs []f
}

// Enabled implements slog.Handler.
func (l *handler) Enabled(ctx context.Context, level slog.Level) bool {
	return l.h.Enabled(ctx, level)
}

// Handle implements slog.Handler.
func (l *handler) Handle(ctx context.Context, r slog.Record) error {
	for _, f := range l.fs {
		r.AddAttrs(f(ctx))
	}
	return l.h.Handle(ctx, r)
}

// WithAttrs implements slog.Handler.
func (l *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return l.h.WithAttrs(attrs)
}

// WithGroup implements slog.Handler.
func (l *handler) WithGroup(name string) slog.Handler {
	return l.h.WithGroup(name)
}

func Info(ctx context.Context, msg string, args ...any) {
	log(ctx, slog.LevelInfo, msg, args...)
}

func log(ctx context.Context, lv slog.Level, msg string, args ...any) {
	l := slog.Default()
	if !l.Enabled(ctx, lv) {
		return
	}

	// see: https://pkg.go.dev/log/slog#example-package-Wrapping
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip [Callers, log, info]

	r := slog.NewRecord(time.Now(), lv, msg, pcs[0])
	r.Add(args...)

	_ = l.Handler().Handle(ctx, r)
}
