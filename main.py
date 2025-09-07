import re
import glob
import os
from typing import Set
import json

def normalize_card_name(card_name):
    """Cleans card names for matching."""
    # TODO: Swap to normalization deck codes use?
    if not isinstance(card_name, str): return None
    name = re.sub(r'^\(\d+\)\s*', '', card_name)
    name = name.replace('-', ' ')
    name = re.sub(r'[^a-zA-Z0-9\s]', '', name)
    return name.lower().strip()

# TODO: Move
definitions = []
normalized_definitions = []
with open("rules.json", "r") as file:
    rules = json.load(file)
    for rule in rules["definitions"]:
        norm_rule = {k: v for k, v in rule.items()}
        for key in ['core_cards', 'at_least_one_of', 'banned_cards']:
            if key in norm_rule:
                norm_rule[key] = {normalize_card_name(c) for c in norm_rule[key]}
        normalized_definitions.append(norm_rule)


def generate_rule_explanation(rule):
    """Creates a human-readable explanation of a rule."""
    parts = []
    if rule.get('core_cards'):
        parts.append(f"Must contain ALL of: [{', '.join(rule['core_cards'])}]")
    if rule.get('at_least_one_of'):
        parts.append(f"Must contain AT LEAST ONE of: [{', '.join(rule['at_least_one_of'])}]")
    if rule.get('banned_cards'):
        parts.append(f"Must NOT contain: [{', '.join(rule['banned_cards'])}]")
    return "; ".join(parts)

def find_most_recent_file(pattern):
    """Finds the most recently modified file in the current directory that matches a pattern."""
    list_of_files = glob.glob(pattern)
    return max(list_of_files, key=os.path.getmtime) if list_of_files else None


def get_archetype_from_deckcode(deck: str, baseEncoded=True) -> str:
    # TODO: Implement
    raise Exception("Not implemented")


def get_archetype(deck: Set[str]) -> str:
    """ 
    Returns a string describing a deck's archetype

    deck - Set of normalized names of cards in a deck (Ex: {"thor", "beta ray bill", "jane foster mighty thor"})
    """
    for rule in normalized_definitions:
        core = rule.get('core_cards', set())
        at_least_one = rule.get('at_least_one_of', set())
        banned = rule.get('banned_cards', set())
        if core.issubset(deck) and (not at_least_one or at_least_one.intersection(deck)) and (not banned or not banned.intersection(deck)):
            return rule['name']
    return 'Miscellaneous / Other'

def main():
    test_deck = set(["agatha harkness", "hela", "red shift"])
    # TODO: Set up unit tests
    print(get_archetype(test_deck))
    print(get_archetype(set(["thor", "beta ray bill", "jane foster mighty thor"])))

if __name__ == '__main__':
    main()
