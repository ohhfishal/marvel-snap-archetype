from typing import Set, Tuple
import base64
import json
import logging
import string

if __name__ == "__main__":
    logging.basicConfig(level=logging.DEBUG)

logger = logging.getLogger(__name__)

RULE_TYPES = ["core_cards", "at_least_one_of", "banned_cards"]

_NORMALIZE_CACHE = {
    "Mister Negative": "MrNgtvA",
    "The First Ghost Rider": "GhstThFrstRdr12",
    "Jane Foster Mighty Thor": "JnFstrA",
    "M.O.D.O.K.": "Mdk5",
    "Sam Wilson Captain America": "SmWlsn9",
    # TODO: Need the codes for cards I don't have...
}


def decode_deckcode(deckcode: str) -> Set[str]:
    """
    Converts the raw Marvel Snap's copied deckcode to a set of Marvel Snap normalized card names (Ex: "SpdrMn9").
    NOTE: Currently only supports the most recent version of deck encoding introduced in 2024.
    """
    lines = [
        line
        for line in deckcode.splitlines()
        if not line.strip().startswith("#") and len(line.strip()) != 0
    ]
    if len(lines) != 1:
        raise ValueError(
            f"invalid deckcode: expected 1 non-commented line got {len(lines)}"
        )
    return {card for card in base64.b64decode(lines[0]).decode("utf-8").split(",")}


def normalize_card_name(name: str) -> str:
    """
    Normalizes a Marvel Snap card name to the internal shortened encoding.
    (Ex: Spider-Man -> "SpdrMn9")
    """
    if normalized := _NORMALIZE_CACHE.get(name, None):
        return normalized

    count = 1
    iterable = name.lstrip("The").lstrip(" ")
    builder = iterable[0]
    for char in iterable[1:]:
        if char in string.punctuation or char.isspace():
            continue
        if char not in "aeiouy":
            builder += char
        count += 1
    hex_count = f"{count:x}".upper()
    normalized = f"{builder}{hex_count}"
    _NORMALIZE_CACHE[name] = normalized
    return normalized


def generate_rule_explanation(rule):
    """Creates a human-readable explanation of a rule."""
    parts = []
    if cards := rule.get("core_cards"):
        parts.append(f"Must contain ALL of: [{', '.join(cards)}]")
    if cards := rule.get("at_least_one_of"):
        parts.append(f"Must contain AT LEAST ONE of: [{', '.join(cards)}]")
    if cards := rule.get("banned_cards"):
        parts.append(f"Must NOT contain: [{', '.join(cards)}]")
    return "; ".join(parts)


# TODO: Move
definitions = []
archetypes = set()
# TODO: Probably doesn't work when imported by another file
with open("rules.json", "r") as file:
    rules = json.load(file)
    for rule in rules["definitions"]:
        for key, value in rule.items():
            if key in RULE_TYPES:
                rule[key] = {normalize_card_name(c) for c in value}
        definitions.append(rule)
        archetypes.add(rule.get("archetype", rule["name"]))

logger.debug(f"Loaded {len(definitions)} rules")
logger.debug(f"Supported Archetypes: {archetypes}")


def get_archetype_from_deckcode(deck: str) -> str:
    # TODO: Write tests to capture this
    return get_archetype(decode_deckcode(deck))


def get_archetype(deck: Set[str]) -> Tuple[str, str]:
    """
    Returns a tuple contain a decks name and archetype

    deck - Set of normalized names of cards in a deck (Ex: {"Hd4", "Zb4", "SmWlsn9"})
    """
    for rule in definitions:
        core = rule.get("core_cards", set())
        at_least_one = rule.get("at_least_one_of", set())
        banned = rule.get("banned_cards", set())
        if (
            core.issubset(deck)
            and (not at_least_one or at_least_one.intersection(deck))
            and (not banned or not banned.intersection(deck))
        ):
            return rule["name"], rule.get("archetype", rule["name"])
    return "Miscellaneous / Other", "Other"
