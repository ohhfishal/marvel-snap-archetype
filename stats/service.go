package stats

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
)

var tournaments = map[string]string{
	"marvel-snap-golden-gauntlet-world-championship-qualifier-1": "Golden Gauntlet World Championship Qualifier 1",
	"marvel-snap-golden-gauntlet-world-championship-qualifier-2": "Golden Gauntlet World Championship Qualifier 2",
	"marvel-snap-golden-gauntlet-world-championship-qualifier-3": "Golden Gauntlet World Championship Qualifier 3",
}

type ServiceOptions struct {
	ApiKey string `env:"API_KEY" help:"TopDeck API Key"`
}

type Service struct {
	apiKey    string
	outputDir string
	logger    *slog.Logger
}

func NewService(logger *slog.Logger, opts ServiceOptions) (*Service, error) {
	service := Service{}

	service.apiKey = opts.ApiKey
	if service.apiKey == "" {
		return nil, errors.New("missing required field: ApiKey")
	}

	outputDir, err := os.MkdirTemp(os.TempDir(), "snap")
	if err != nil {
		return nil, fmt.Errorf("creating data directory: %w", err)
	}
	service.outputDir = outputDir
	logger.Info("created data directory", "dir", service.outputDir)

	return &service, nil
}

type Archetype struct {
	Name       string
	Derivation string
	Count      int
}

func (service *Service) GetArchetypes(ctx context.Context, tid string) ([]Archetype, error) {
	if err := service.refreshStats(ctx, tid); err != nil {
		return nil, fmt.Errorf("refreshing stats: %w", err)
	}

	// TODO: Parse the file

	return nil, nil
}

type deckEntry struct {
	Category   string
	Derivation string
	Count      int
	Count64    int
	Count32    int
	Count16    int
}

func (service *Service) refreshStats(ctx context.Context, tid string) error {
	if _, ok := tournaments[tid]; !ok {
		return fmt.Errorf("invalid tid: %s", tid)
	}

	outputDir := path.Join(service.outputDir, tid)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed creating output directory: %w", err)
	}

	cmd := Analysis{
		TID:    tid,
		ApiKey: service.apiKey,
		Output: outputDir,
	}

	if err := cmd.Run(ctx, service.logger); err != nil {
		return fmt.Errorf("getting stats: %w", err)
	}
	return nil
}

func (service *Service) Tournaments() map[string]string {
	return tournaments
}
