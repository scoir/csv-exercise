from unittest import TestCase

from src.models import Person


class TestPerson(TestCase):
    def test_new_person_invalid_params(self):
        with self.assertRaises(ValueError):
            Person(internal_id=None)
        with self.assertRaises(ValueError):
            Person(internal_id="42")
        with self.assertRaises(ValueError):
            Person(internal_id="42", first_name="Bobby")
        with self.assertRaises(ValueError):
            Person(internal_id="42", first_name="Bobby", last_name="tables")
