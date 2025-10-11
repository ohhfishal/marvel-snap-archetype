package job

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ohhfishal/marvel-snap-archetype/assets"
	"log/slog"
	"os"
	"strings"
)

var (
	definitions []Rule
)

type Rule struct {
	Name         string   `json:"name"`
	Archetype    *string  `json:"archetype,omitempty"`
	CoreCards    []string `json:"core_cards,omitempty"`
	AtLeastOneOf []string `json:"at_least_one_of,omitempty"`
	BannedCards  []string `json:"banned_cards,omitempty"`
}

type RulesData struct {
	Definitions []Rule `json:"definitions"`
}

func init() {
	var rulesData RulesData
	if err := json.Unmarshal(assets.RulesJSON, &rulesData); err != nil {
		panic(fmt.Sprintf("failed to unmarshal rules JSON: %v", err))
	}
	definitions = rulesData.Definitions
}

func DeckStats(logger *slog.Logger, path string, standings []Standing, cuts []int) error {
	if len(cuts) == 0 {
		return errors.New("must include at least one cut of players")
	}

	// TODO: This is janky code

	// TopCut: Archetype: Variants: count
	mappings := map[int]map[string]map[string]int{}
	for _, cut := range cuts {
		mappings[cut] = map[string]map[string]int{}
	}

	for _, player := range standings {
		for topcut, mapping := range mappings {
			if player.Standing > topcut {
				continue
			}

			name, archetype := GetArchetype(player.Deck.Cards)

			if strings.Contains(name, "Misc") {
				slog.Warn("deck classified as other", "deck", player.Deck.Cards)
			}

			// Check if the high level archetype is there
			if _, ok := mapping[archetype]; !ok {
				mapping[archetype] = map[string]int{}
			}

			// Increment the counter
			if _, ok := mapping[archetype][name]; !ok {
				mapping[archetype][name] = 0
			}
			mapping[archetype][name] += 1
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

	header := []string{"Category", "Derivation", "Count"}
	for _, cut := range cuts[1:] {
		header = append(header, fmt.Sprintf("Top %d", cut))
	}

	if err := writer.Write(header); err != nil {
		return fmt.Errorf("writing to file: %w", err)
	}
	// mappings := map[int]map[string]map[string]int{}

	for archetype, subtypes := range mappings[cuts[0]] {
		records := [][]string{}
		totals := map[int]int{cuts[0]: 0}
		for name, count := range subtypes {
			record := []string{
				archetype,
				name,
				fmt.Sprintf("%d", count),
			}
			totals[cuts[0]] += count
			for _, cut := range cuts[1:] {
				if _, ok := totals[cut]; !ok {
					totals[cut] = 0
				}

				if count, ok := mappings[cut][archetype][name]; ok {
					record = append(record, fmt.Sprintf("%d", count))
					totals[cut] += count
				} else {
					record = append(record, "0")
				}

			}
			records = append(records, record)
		}

		totalRecord := []string{
			archetype,
			"Total",
			fmt.Sprintf("%d", totals[cuts[0]]),
		}
		slog.Debug("sum", "archetype", archetype, "total", totalRecord)

		for _, cut := range cuts[1:] {
			totalRecord = append(totalRecord, fmt.Sprintf("%d", totals[cut]))
		}

		if err := writer.WriteAll([][]string{totalRecord}); err != nil {
			return fmt.Errorf("writing to file: %w", err)
		}

		if err := writer.WriteAll(records); err != nil {
			return fmt.Errorf("writing to file: %w", err)
		}
	}
	return nil
}

func GetArchetype(deck map[string]any) (string, string) {
	for _, rule := range definitions {
		coreMatch := isSubset(rule.CoreCards, deck)
		atLeastOneMatch := len(rule.AtLeastOneOf) == 0 || hasIntersection(rule.AtLeastOneOf, deck)
		bannedMatch := len(rule.BannedCards) == 0 || !hasIntersection(rule.BannedCards, deck)

		if coreMatch && atLeastOneMatch && bannedMatch {
			archetype := rule.Name
			if rule.Archetype != nil {
				archetype = *rule.Archetype
			}
			return rule.Name, archetype
		}
	}
	return "Miscellaneous / Other", "Other"
}

func isSubset(subset []string, superset map[string]any) bool {
	for _, key := range subset {
		if _, ok := superset[key]; !ok {
			return false
		}
	}
	return true
}

func hasIntersection(set1 []string, set2 map[string]any) bool {
	for _, key := range set1 {
		if _, exists := set2[key]; exists {
			return true
		}
	}
	return false
}

// TODO: Remove if TopDeck API returns DeckObj during GG3
// func GetArchetypeFromDeckcode(deckcode string) (string, string, error) {
// 	deck, err := decodeDeckcode(deckcode)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	name, category := GetArchetype(deck)
// 	return name, category, nil
// }

// func decodeDeckcode(deckcode string) (map[string]any, error) {
// 	scanner := bufio.NewScanner(strings.NewReader(strings.ReplaceAll(deckcode, "\\n", "\n")))
// 	lines := []string{}
//
// 	for scanner.Scan() {
// 		trimmed := strings.TrimSpace(scanner.Text())
// 		if !strings.HasPrefix(trimmed, "#") && len(trimmed) != 0 {
// 			lines = append(lines, trimmed)
// 		}
// 	}
//
// 	if len(lines) != 1 {
// 		return nil, fmt.Errorf("invalid deckcode: expected 1 non-commented line got %d", len(lines))
// 	}
//
// 	// This fixes some issues from people (1) messing with the deckcode, but is the most I'll fix
// 	var line = lines[0]
// 	for range len(line) % 4 {
// 		line = line + "="
// 	}
//
// 	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(line))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to decode base64: %w", err)
// 	}
//
// 	cards := strings.Split(string(decoded), ",")
// 	deck := make(map[string]any)
// 	for _, card := range cards {
// 		deck[card] = nil
// 	}
//
// 	return deck, nil
// }
