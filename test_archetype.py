import unittest
from main import get_archetype


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
        for i, (expected, input_deck) in enumerate(self.test_cases):
            with self.subTest(case=i, expected=expected):
                result = get_archetype(input_deck)
                self.assertEqual(result, expected)


if __name__ == "__main__":
    unittest.main()
