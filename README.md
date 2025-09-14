# Marvel Snap Deck Analysis

Initial schema based heavily on F0x's work converting GG# decks to archetypes

The CLI uses the [TopDeck.gg](https://topdeck.gg) API for getting tournament data.

*Note:* Deck stats ignore players that could not submit a valid deck. Card stats may still count them if TopDeck manually set the decklist.

## Example using the Golang CLI

```bash
go install github.com/ohhfishal/marvel-snap-archetype@main
export API_KEY=<TopDeck GG KEY!>
go run . analysis "marvel-snap-golden-gauntlet-world-championship-qualifier-2"
```
