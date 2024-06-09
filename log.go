package examplegoslog

import (
	"context"
	"log/slog"
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
