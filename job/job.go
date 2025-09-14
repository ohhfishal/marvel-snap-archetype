package job

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Standing struct {
	Name     string `json:"name"`
	ID       string `json:"id"`
	Standing int    `json:"standing"`
	Decklist string `json:"decklist"` // Marvel Snap deck code
	Deck     struct {
		Cards map[string]any `json:"Decklist"` // Ex: "Kraven": { "id": "#", "count": 1 }
	} `json:"deckObj"`
}

type TournamentResponse struct {
	TID       string     `json:"TID"`
	Standings []Standing `json:"standings"`
}

func GetResults(ctx context.Context, api_key string, tid string) (*TournamentResponse, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://topdeck.gg/api/v2/tournaments/%s", tid),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	request.Header.Add("Authorization", api_key)

	client := &http.Client{}

	rawResponse, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer rawResponse.Body.Close()

	if rawResponse.StatusCode >= 400 || rawResponse.StatusCode < 200 {
		return nil, fmt.Errorf("request failed: %d", rawResponse.StatusCode)
	}

	var response TournamentResponse
	if err := json.NewDecoder(rawResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("got invalid json: %w", err)
	}
	return &response, nil
}
