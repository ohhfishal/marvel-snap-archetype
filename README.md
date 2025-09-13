# Marvel Snap Deck Analysis

Initial schema based heavily on F0x's work converting GG# decks to archetypes

## Example using the Golang CLI

The CLI uses the [TopDeck.gg](https://topdeck.gg) API for getting tournament data.

```bash
go install github.com/ohhfishal/marvel-snap-archetype@main
export API_KEY=<TopDeck GG KEY!>
go run . analysis "marvel-snap-golden-gauntlet-world-championship-qualifier-2"
```
