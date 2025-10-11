package job

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log/slog"
	"os"
)

func CardStats(path string, standings []Standing, cuts []int) error {
	if len(cuts) == 0 {
		return errors.New("must include at least one cut of players")
	}

	mappings := map[int]map[string]int{}
	for _, cut := range cuts {
		mappings[cut] = map[string]int{}
	}

	slog.Info("parsing cards", "num_players", len(standings))
	for _, player := range standings {
		for topcut, mapping := range mappings {
			if player.Standing > topcut {
				continue
			}
			if len(player.Deck.Cards) != 12 {
				slog.Warn("Deck not made of 12 cards", "player", player)
			}
			slog.Debug("player using", "deck", player.Deck.Cards)
			for card, _ := range player.Deck.Cards {
				if _, ok := mapping[card]; !ok {
					mapping[card] = 0
				}
				mapping[card] += 1
			}
		}
	}

	file, err := os.OpenFile(
		path,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0644,
	)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Name", "Count"}
	for _, cut := range cuts[1:] {
		header = append(header, fmt.Sprintf("Top %d", cut))
	}

	if err := writer.Write(header); err != nil {
		return fmt.Errorf("writing to file: %w", err)
	}

	for card, count := range mappings[cuts[0]] {
		record := []string{
			card,
			fmt.Sprintf("%d", count),
		}

		for _, cut := range cuts[1:] {
			if count, ok := mappings[cut][card]; ok {
				record = append(record, fmt.Sprintf("%d", count))
			} else {
				record = append(record, "0")
			}
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("writing to file: %w", err)
		}
	}
	return nil
}
