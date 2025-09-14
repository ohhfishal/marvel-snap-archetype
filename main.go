package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/ohhfishal/marvel-snap-archetype/job"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
)

type CLI struct {
	Analysis Analysis `cmd:"" help:"Get stats on a tournament."`
}

type Analysis struct {
	TID    string `arg:"" help:"Tournament ID (Check TopDeck URL)."`
	ApiKey string `env:"API_KEY" help:"TopDeck API Key"`
	Output string `default:"output" type:"path" help:"Output directory path"`
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
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	if err := kongCtx.Run(logger); err != nil {
		logger.Error("failed", slog.Any("error", err))
		os.Exit(1)
	}
}

func (config *Analysis) Run(ctx context.Context, logger *slog.Logger) error {
	if config.ApiKey == "" {
		return errors.New("missing env: API_KEY")
	}

	if err := os.MkdirAll(config.Output, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	response, err := job.GetResults(ctx, config.ApiKey, config.TID)
	if err != nil {
		return fmt.Errorf("getting tournament results: %w", err)
	}

	var wg sync.WaitGroup
	wg.Go(func() {
		if err := job.CardStats(
			filepath.Join(config.Output, "cards.csv"),
			response.Standings,
			[]int{1024, 64, 32, 16}, // TODO: Make configurable and passed in somehow
		); err != nil {
			logger.Error("job failed",
				slog.Any("error", err),
				slog.String("job", "cards"),
			)
		}
	})
	wg.Wait()
	return nil
}
