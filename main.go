package main

import (
	"context"
	"github.com/alecthomas/kong"
	"github.com/ohhfishal/marvel-snap-archetype/server"
	"github.com/ohhfishal/marvel-snap-archetype/stats"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type CLI struct {
	Analysis stats.Analysis `cmd:"" help:"Get stats on a tournament."`
	Serve    server.CMD     `cmd:"" help:"Run stats hosting server."`
}

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	var cli CLI
	kongCtx := kong.Parse(
		&cli,
		kong.BindTo(ctx, new(context.Context)),
	)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := kongCtx.Run(logger); err != nil {
		logger.Error("failed", slog.Any("error", err))
		os.Exit(1)
	}
}
