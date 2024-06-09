package main

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"os"

	examplegoslog "github.com/blck-snwmn/example-go-slog"
)

type key struct{}

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, key{}, "ID-xxxxx")

	h := examplegoslog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
		func(ctx context.Context) slog.Attr {
			fmt.Println("--called")
			return slog.Any("id", cmp.Or(ctx.Value(key{}), "none"))
		},
	)
	logger := slog.New(h)

	logger.Info("do1")
	logger.InfoContext(ctx, "do2")
	logger.InfoContext(ctx, "do3")
	logger.Info("do4")
	logger.InfoContext(ctx, "do5")
	logger.Debug("do5")
	logger.DebugContext(ctx, "do6")
	logger.InfoContext(ctx, "do7")
}
