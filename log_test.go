package examplegoslog

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strconv"
	"testing"
)

type c struct {
	count int
}

func (c *c) String() string {
	c.count++
	return strconv.Itoa(c.count)
}

func TestSlog(t *testing.T) {
	t.Run("Count up each time slog.Info is called.", func(t *testing.T) {
		c := &c{}
		slog.Info("xxx", "count", c)
		slog.Info("xxx", "count", c)

		if c.count != 2 {
			t.Errorf("c.count = %d; want 2", c.count)
		}
	})
	t.Run("Counts do not increase if no info log is output", func(t *testing.T) {
		h := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		})
		logger := slog.New(h)

		c := &c{}

		logger.Info("xxx", "count", c)
		logger.Info("xxx", "count", c)

		if c.count != 0 {
			t.Errorf("c.count = %d; want 0", c.count)
		}
	})
	t.Run("If you use `with` to add an attribute, it will only count up once.", func(t *testing.T) {
		h := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		logger := slog.New(h)
		c := &c{}
		logger = logger.With("count", c)

		if c.count != 1 {
			t.Errorf("c.count = %d; want 1", c.count)
		}

		logger.Info("xxx")
		logger.Info("xxx")

		if c.count != 1 {
			t.Errorf("c.count = %d; want 1", c.count)
		}
	})
}

func TestHandler(t *testing.T) {
	t.Run("Count up each time slog.Info is called.", func(t *testing.T) {
		c := &c{}
		h := New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}), func(ctx context.Context) slog.Attr {
			return slog.Any("count", c)
		})
		logger := slog.New(h)

		logger.Info("xxx")
		logger.Info("xxx")

		if c.count != 2 {
			t.Errorf("c.count = %d; want 2", c.count)
		}
	})

	t.Run("Counts do not increase if no info log is output", func(t *testing.T) {
		c := &c{}
		h := New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		}), func(ctx context.Context) slog.Attr {
			return slog.Any("count", c)
		})
		logger := slog.New(h)

		logger.Info("xxx")
		logger.Info("xxx")

		if c.count != 0 {
			t.Errorf("c.count = %d; want 0", c.count)
		}
	})
}
