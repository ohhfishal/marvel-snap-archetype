package stats

import (
	"context"
	"errors"
	"fmt"
	"github.com/ohhfishal/marvel-snap-archetype/job"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

type Analysis struct {
	TID    string `arg:"" help:"Tournament ID (Check TopDeck URL)."`
	ApiKey string `env:"API_KEY" help:"TopDeck API Key"`
	Output string `default:"output" type:"path" help:"Output directory path"`
}

func (config *Analysis) Run(ctx context.Context, logger *slog.Logger) error {
	if config.ApiKey == "" {
		return errors.New("missing env: API_KEY")
	}

	if logger == nil {
		logger = slog.New(slog.DiscardHandler)
	}

	if err := os.MkdirAll(config.Output, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	response, err := job.GetResults(ctx, config.ApiKey, config.TID)
	if err != nil {
		return fmt.Errorf("getting tournament results: %w", err)
	}

	cuts := []int{1024, 64, 32, 16} // TODO: Make configurable

	var wg sync.WaitGroup
	wg.Go(func() {
		if err := job.CardStats(
			filepath.Join(config.Output, "cards.csv"),
			response.Standings,
			cuts,
		); err != nil {
			logger.Error("job failed",
				slog.Any("error", err),
				slog.String("job", "cards"),
			)
		}
	})
	wg.Go(func() {
		if err := job.DeckStats(
			logger.With("job", "decks"),
			filepath.Join(config.Output, "decks.csv"),
			response.Standings,
			cuts,
		); err != nil {
			logger.Error("job failed",
				slog.Any("error", err),
				slog.String("job", "decks"),
			)
		}
	})
	wg.Wait()
	return nil
}
