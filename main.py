import re
from typing import Set, Tuple
import json
import logging

if __name__ == "__main__":
    logging.basicConfig(level=logging.DEBUG)


logger = logging.getLogger(__name__)

RULE_TYPES = ["core_cards", "at_least_one_of", "banned_cards"]


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


def normalize_card_name(card_name):
    """Cleans card names for matching."""
    # TODO: Swap to normalization deck codes use?
    if not isinstance(card_name, str):
        return None
    name = re.sub(r"^\(\d+\)\s*", "", card_name)
    name = name.replace("-", " ")
    name = re.sub(r"[^a-zA-Z0-9\s]", "", name)
    return name.lower().strip()


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


def get_archetype_from_deckcode(deck: str, baseEncoded=True) -> str:
    # TODO: Implement
    raise Exception("Not implemented")


def get_archetype(deck: Set[str]) -> Tuple[str, str]:
    """
    Returns a tuple contain a decks name and archetype

    deck - Set of normalized names of cards in a deck (Ex: {"thor", "beta ray bill", "jane foster mighty thor"})
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
