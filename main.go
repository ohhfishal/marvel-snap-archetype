package main

import (
	"context"
	"encoding/csv"
	json "encoding/json/v2"
	"errors"
	"fmt"
	"github.com/alecthomas/kong"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type CLI struct {
	Analysis Analysis `cmd:"" help:"Get stats on a tournament."`
}

type Analysis struct {
	TID    string `arg:"" help:"Tournament ID (Check TopDeck URL)."`
	ApiKey string `env:"API_KEY" help:"TopDeck API Key"`
	Output string `default:"out.csv" type:"path" help:"Output file path"`
}

type TournamentResponse struct {
	TID       string `json:"TID"`
	Standings []struct {
		Name     string `json:"name"`
		Decklist string `json:"decklist"` // Marvel Snap deck code
		Deck     struct {
			Cards map[string]any `json:"Decklist"` // Ex: "Kraven": { "id": "#", "count": 1 }
		} `json:"deckObj"`
	} `json:"standings"`
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

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://topdeck.gg/api/v2/tournaments/%s", config.TID),
		nil,
	)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	request.Header.Add("Authorization", config.ApiKey)

	client := &http.Client{}

	rawResponse, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("making request: %w", err)
	}
	defer rawResponse.Body.Close()

	if rawResponse.StatusCode >= 400 {
		return fmt.Errorf("request failed: %d", rawResponse.StatusCode)
	}

	var response TournamentResponse
	if err := json.UnmarshalRead(rawResponse.Body, &response); err != nil {
		return fmt.Errorf("got invalid json: %w", err)
	}

	mapping := map[string]int{}
	for _, player := range response.Standings {
		for card, _ := range player.Deck.Cards {
			if _, ok := mapping[card]; !ok {
				mapping[card] = 0
			}
			mapping[card] += 1
		}
	}

	file, err := os.OpenFile(config.Output, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"Name", "Count"}); err != nil {
		return fmt.Errorf("writing to file: %w", err)
	}

	for card, count := range mapping {
		if err := writer.Write([]string{card, fmt.Sprintf("%d", count)}); err != nil {
			return fmt.Errorf("writing to file: %w", err)
		}
	}
	return nil
}
