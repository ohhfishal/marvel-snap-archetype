import unittest
from main import get_archetype, decode_deckcode, normalize_card_name


class TestNormalization(unittest.TestCase):
    test_cases = [
        ("The Hood", "Hd4"),
        ("Zabu", "Zb4"),
        ("Sam Wilson Captain America", "SmWlsn9"),
        ("Spider-Man", "SpdrMn9"),
        ("Symbiote Spider-Man ", "SmbtSpdrMn11"),
        ("Agatha Harkness", "AgthHrknssE"),
        ("X-23", "X233"),
        ("M.O.D.O.K.", "Mdk5"),
        ("Human Torch First Steps", "HmnTrchFrstStps14"),
        ("Spider-Woman", "SpdrWmnB"),
        ("Spider-Man 2099", "SpdrMn2099D"),
        ("Mockingbird", "MckngbrdB"),
        ("Mister Negative", "MrNgtvA"),
        ("Agent 13", "Agnt137"),
        ("The First Ghost Rider", "GhstThFrstRdr12"),
        ("Professor X", "PrfssrXA"),
        ("Jane Foster Mighty Thor", "JnFstrA"),
        ("Agent Venom", "AgntVnmA"),
        ("Invisible Woman First Steps", "InvsblWmnFrstStps18"),
        ("Phastos", "Phsts7"),
        ("Mobius M. Mobius", "MbsMMbsD"),
        ("Negasonic Teenage Warhead", "NgsncTngWrhd17"),
        ("Brood", "Brd5"),
        ("U.S. Agent", "USAgnt7"),
    ]

    def test_normalization(self):
        for i, (human_readable, encoded) in enumerate(self.test_cases):
            with self.subTest(case=i, expected=human_readable, encoded=encoded):
                result = normalize_card_name(human_readable)
                self.assertEqual(result, encoded)


class TestDecodeDeckcode(unittest.TestCase):
    test_cases = [
        (
            set(
                [
                    "Hd4",
                    "Zb4",
                    "SmWlsn9",
                    "SpdrMn9",
                    "SmbtSpdrMn11",
                    "AgthHrknssE",
                    "X233",
                    "Mdk5",
                    "HmnTrchFrstStps14",
                    "SpdrWmnB",
                    "SpdrMn2099D",
                    "MckngbrdB",
                ]
            ),
            """
            # (1) The Hood
            # (1) Zabu
            # (1) X-23
            # (2) Sam Wilson Captain America
            # (2) Spider-Man
            # (3) Human Torch First Steps
            # (4) Symbiote Spider-Man
            # (5) M.O.D.O.K.
            # (5) Spider-Woman
            # (5) Spider-Man 2099
            # (6) Mockingbird
            # (6) Agatha Harkness
            #
            SGQ0LFpiNCxTbVdsc245LFNwZHJNbjksU21idFNwZHJNbjExLEFndGhIcmtuc3NFLFgyMzMsTWRrNSxIbW5UcmNoRnJzdFN0cHMxNCxTcGRyV21uQixTcGRyTW4yMDk5RCxNY2tuZ2JyZEI=
            #
            # To use this deck, copy it to your clipboard and paste it from the deck editing menu in MARVEL SNAP.
            """,
        ),
        (
            set(
                [
                    "MrNgtvA",
                    "Agnt137",
                    "GhstThFrstRdr12",
                    "PrfssrXA",
                    "JnFstrA",
                    "AgntVnmA",
                    "InvsblWmnFrstStps18",
                    "Phsts7",
                    "MbsMMbsD",
                    "NgsncTngWrhd17",
                    "Brd5",
                    "USAgnt7",
                ]
            ),
            """
            # (1) Agent 13
            # (2) U.S. Agent
            # (2) Agent Venom
            # (2) The First Ghost Rider
            # (3) Brood
            # (3) Negasonic Teenage Warhead
            # (3) Mobius M. Mobius
            # (3) Phastos
            # (3) Invisible Woman First Steps
            # (4) Mister Negative
            # (5) Professor X
            # (5) Jane Foster Mighty Thor
            #
            TXJOZ3R2QSxBZ250MTM3LEdoc3RUaEZyc3RSZHIxMixQcmZzc3JYQSxKbkZzdHJBLEFnbnRWbm1BLEludnNibFdtbkZyc3RTdHBzMTgsUGhzdHM3LE1ic01NYnNELE5nc25jVG5nV3JoZDE3LEJyZDUsVVNBZ250Nw==
            #
            # To use this deck, copy it to your clipboard and paste it from the deck editing menu in MARVEL SNAP.
         """,
        ),
    ]

    def test_parsing(self):
        for i, (expected, input_deck) in enumerate(self.test_cases):
            with self.subTest(case=i, expected=expected):
                result = decode_deckcode(input_deck)
                self.assertEqual(result, expected)


class TestGetArchetype(unittest.TestCase):
    test_cases = [
        (
            ("Agatha Hela", "Agatha"),
            {"agatha harkness", "hela", "red shift"},
        ),
        (
            ("Thors", "Thors"),
            {"thor", "beta ray bill", "jane foster mighty thor"},
        ),
        (
            ("Cerebro-2", "Cerebro"),
            {"cerebro", "mystique", "lasher"},
        ),
        (
            ("Cerebro-3", "Cerebro"),
            {"cerebro", "mystique", "scarlet witch"},
        ),
        # TODO: Add more
    ]

    def test_archetypes(self):
        return
        for i, (expected, input_deck) in enumerate(self.test_cases):
            with self.subTest(case=i, expected=expected):
                result = get_archetype(input_deck)
                self.assertEqual(expected, result)


if __name__ == "__main__":
    unittest.main()
