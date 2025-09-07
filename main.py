import re
import glob
import os
from typing import Set

ARCHETYPE_DEFINITIONS = [
    {'name': 'Agatha Hela', 'parent': 'Agatha', 'core_cards': ['Agatha Harkness', 'Hela', 'Red Shift']},
    {'name': 'Agatha Redshift', 'parent': 'Agatha', 'core_cards': ['Agatha Harkness', 'Red Shift']},
    {'name': 'Non-Agatha Hela', 'core_cards': ['Hela'], 'banned_cards': ['Agatha Harkness']},
    {'name': 'Apocalypse Discard', 'parent': 'Discard', 'core_cards': ['Dracula', 'M.O.D.O.K.', 'Apocalypse', 'Khonshu']},
    {'name': 'Bullseye Discard', 'parent': 'Discard', 'core_cards': ['M.O.D.O.K.', 'Bullseye', 'Swarm', 'Daken']},
    {'name': 'Wiccan Domino', 'parent': 'Wiccan', 'core_cards': ['Quicksilver', 'Wiccan', 'Domino']},
    {'name': 'Wiccan', 'core_cards': ['Quicksilver', 'Wiccan']},
    {'name': 'Move Bounce', 'parent': 'Move', 'core_cards': ['Human Torch', 'Madame Web', 'Beast', 'Toxin']},
    {'name': 'Move', 'core_cards': ['Human Torch', 'Madame Web']},
    {'name': 'Destroy V.Hand', 'parent': 'Victoria Hand', 'core_cards': ['Moira X', 'The Hood', 'Victoria Hand', 'Frigga']},
    {'name': 'Moon Girl V.Hand', 'parent': 'Victoria Hand', 'core_cards': ['Victoria Hand', 'Frigga', 'Moon Girl', 'Misery', 'The Hood'], 'banned_cards': ['Werewolf By Night']},
    {'name': 'Werewolf V.Hand', 'parent': 'Victoria Hand', 'core_cards': ['Victoria Hand', 'Frigga', 'Misery', 'The Hood', 'Werewolf By Night']},
    {'name': 'Other V.Hand', 'parent': 'Victoria Hand', 'core_cards': ['Victoria Hand']},
    {'name': 'Cerebro-2', 'parent': 'Cerebro', 'core_cards': ['Cerebro', 'Mystique', 'Mister Sinister', 'Brood']},
    {'name': 'Cerebro-3', 'parent': 'Cerebro', 'core_cards': ['Cerebro', 'Mystique'], 'at_least_one_of': ['Bast', 'Scarlet Witch', 'Cosmo', 'Negasonic Teenage Warhead']},
    {'name': 'Cerebro-4', 'parent': 'Cerebro', 'core_cards': ['Cerebro', 'Prodigy']},
    {'name': 'Cerebro-5', 'parent': 'Cerebro', 'core_cards': ['Cerebro', 'The Ancient One']},
    {'name': 'Cerebro-6', 'parent': 'Cerebro', 'core_cards': ['Cerebro', 'Moonstone']},
    {'name': 'Thanos Ongoing', 'parent': 'Thanos', 'core_cards': ['Thanos', 'Spectrum']},
    {'name': 'Arishem Thanos', 'parent': 'Arishem', 'core_cards': ['Arishem', 'Thanos']},
    {'name': 'Thanos', 'core_cards': ['Thanos']},
    {'name': 'Arishem', 'core_cards': ['Arishem']},
    {'name': 'Agamotto', 'core_cards': ['Agamotto'], 'banned_cards': ['Thanos', 'Arishem']},
    {'name': 'Doom 2099 EoT', 'parent': 'End of Turn', 'core_cards': ['Invisible Woman First Steps', 'Doctor Doom 2099']},
    {'name': 'Invisible Woman Copy', 'parent': 'End of Turn', 'core_cards': ['Invisible Woman First Steps', 'Prodigy']},
    {'name': 'End of Turn', 'core_cards': ['Invisible Woman First Steps']},
    {'name': 'Sera Hitmonkey', 'parent': 'Sera', 'core_cards': ['Sera', 'Hit-Monkey'], 'at_least_one_of': ['Bishop', 'Mysterio']},
    {'name': 'Sera Control', 'parent': 'Sera', 'core_cards': ['Sera', 'Shang-Chi']},
    {'name': 'Galactus Ramp', 'parent': 'Electro Ramp', 'core_cards': ['Electro', 'Blink', 'Galactus']},
    {'name': 'Electro Ramp', 'core_cards': ['Electro', 'Blink']},
    {'name': 'Silver Surfer', 'parent': 'Brood', 'core_cards': ['Silver Surfer', 'Brood']},
    {'name': 'Handbuff', 'parent': 'Brood', 'core_cards': ['Gwenpool', 'Brood'], 'banned_cards': ['Silver Surfer']},
    {'name': 'Tribunal', 'core_cards': ['The Living Tribunal', 'Iron Man', 'Onslaught']},
    {'name': 'Werewolf Pile', 'core_cards': ['Merlin', 'Werewolf By Night', 'Misery', 'The Hood']},
    {'name': 'Mr. Negative', 'core_cards': ['Mister Negative', 'Gorr the God Butcher']},
    {'name': 'Deadpool Destroy', 'core_cards': ['Deadpool', 'Venom', 'Death']},
    {'name': 'Zoo', 'core_cards': ['Squirrel Girl', 'Marvel Boy'], 'banned_cards': ['Invisible Woman First Steps']},
    {'name': 'Shuri Nimrod', 'core_cards': ['Kid Omega', 'Shuri', 'Nimrod']},
    {'name': 'Squatchbird', 'core_cards': ['Mysterio', 'Sasquatch', 'Mockingbird']},
    {'name': 'Zabu Tech', 'core_cards': ['Zabu', 'Galacta', 'Shang-Chi']},
    {'name': 'Scream', 'core_cards': ['Scream', 'Spider-Man', 'Polaris']},
    {'name': 'Air Walker Morgan', 'core_cards': ['Air Walker', 'Morgan la Fey']},
    {'name': 'Spectrum', 'core_cards': ['Sam Wilson Captain America', 'Captain America', 'Moonstone']},
    {'name': 'Pixie', 'core_cards': ['Pixie']},
    {'name': 'Storm Legion', 'core_cards': ['Storm', 'War Machine', 'Legion']},
    {'name': 'Wong Odin', 'core_cards': ['Wong', 'Odin']},
    {'name': 'Thors', 'core_cards': ['Thor', 'Beta Ray Bill', 'Jane Foster Mighty Thor']},
    {'name': 'Bounce', 'core_cards': ['Rocket Raccoon', 'Toxin', 'Beast']},
    {'name': 'Hydra Stomper', 'core_cards': ['Batroc the Leaper', 'Sam Wilson Captain America', 'Hydra Stomper']},
    {'name': 'Clog', 'core_cards': ['Debrii', 'Green Goblin', 'Titania']},
]

# PARAMETERS
CONVERSION_WIN_THRESHOLD = 5
FILE_PATTERN = 'ggq*-data.xlsx'
OUTPUT_CSV_PATH = 'expert_archetype_report_final.csv'

def normalize_card_name(card_name):
    """Cleans card names for matching."""
    if not isinstance(card_name, str): return None
    name = re.sub(r'^\(\d+\)\s*', '', card_name)
    name = name.replace('-', ' ')
    name = re.sub(r'[^a-zA-Z0-9\s]', '', name)
    return name.lower().strip()

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

# {'name': 'Agatha Hela', 'parent': 'Agatha', 'core_cards': ['Agatha Harkness', 'Hela', 'Red Shift']},
test_deck = set(["agatha harkness", "hela", "red shift"])


# TODO: Move
definitions = ARCHETYPE_DEFINITIONS
normalized_definitions = []
for rule in definitions:
    norm_rule = {k: v for k, v in rule.items()}
    for key in ['core_cards', 'at_least_one_of', 'banned_cards']:
        if key in norm_rule:
            norm_rule[key] = {normalize_card_name(c) for c in norm_rule[key]}
    normalized_definitions.append(norm_rule)

def get_archetype(deck: Set[str]):
    for rule in normalized_definitions:
        core = rule.get('core_cards', set())
        at_least_one = rule.get('at_least_one_of', set())
        banned = rule.get('banned_cards', set())
        if core.issubset(deck) and (not at_least_one or at_least_one.intersection(deck)) and (not banned or not banned.intersection(deck)):
            return rule['name']
    return 'Miscellaneous / Other'

def main():
    print(get_archetype(test_deck))

if __name__ == '__main__':
    main()
