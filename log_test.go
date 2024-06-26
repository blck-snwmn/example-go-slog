package examplegoslog

import (
	"bytes"
	"context"
	"encoding/json"
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

func BenchmarkSlog(b *testing.B) {
	type lkey struct{}

	h := slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(h)
	logger = logger.With("count", "xxx")

	ctx := context.WithValue(context.Background(), lkey{}, logger)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx.Value(lkey{}).(*slog.Logger).Info("xxx")
	}
}

func BenchmarkHandler(b *testing.B) {
	type vkey struct{}

	h := New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}), func(ctx context.Context) slog.Attr {
		v := ctx.Value(vkey{}).(string)
		return slog.String("c", v)
	})
	logger := slog.New(h)

	ctx := context.WithValue(context.Background(), vkey{}, "xxx")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.InfoContext(ctx, "xxx")
	}
}

func TestDo(t *testing.T) {
	var buff bytes.Buffer
	l := slog.New(slog.NewJSONHandler(&buff, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}))
	slog.SetDefault(l)
	Info(context.Background(), "xxx1", "x", "z")

	var m map[string]interface{}
	json.Unmarshal(buff.Bytes(), &m)

	if m["msg"] != "xxx1" {
		t.Errorf("msg = %v; want xxx1", m["msg"])
	}
	if m["level"] != slog.LevelInfo.String() {
		t.Errorf("level = %v; want info", m["level"])
	}
	if m["x"] != "z" {
		t.Errorf("x = %v; want z", m["x"])
	}
	if m["source"] == `{"function":"github.com/blck-snwmn/example-go-slog.TestDo","file":"/home/snowman/dev/github.com/blck-snwmn/example-go-slog/log_test.go","line":146}` {
		t.Errorf("source = %v; want source", m["source"])
	}
}
