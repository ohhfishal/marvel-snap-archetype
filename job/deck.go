package job

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ohhfishal/marvel-snap-archetype/assets"
	"log/slog"
	"os"
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

	mappings := map[int]map[string]int{}
	for _, cut := range cuts {
		mappings[cut] = map[string]int{}
	}

	// TODO: MOVE
	file, err := os.OpenFile(
		path,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0644,
	)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	writer := csv.NewWriter(file)
	defer writer.Flush()
	header := []string{"Name", "Category"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("writing to file: %w", err)
	}

	for _, player := range standings {
		for topcut, _ := range mappings {
			if player.Standing > topcut {
				continue
			}
			name, category := GetArchetype(player.Deck.Cards)
			if err := writer.Write([]string{name, category}); err != nil {
				return fmt.Errorf("writing to file: %w", err)
			}
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
