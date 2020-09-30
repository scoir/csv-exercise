import os
from unittest import TestCase

from src.models import Person, Error
from src.utils import delete_file, write_json, write_csv


class Test(TestCase):
    def setUp(self):
        self.people = [Person(internal_id="42", first_name="Bobby", last_name="tables", phone_number="555-867-5309")]
        self.errors = [Error(row="42", message="Meaning of life not found")]

    def test_delete_file_success(self):
        with open("test.txt", 'w') as f:
            pass

        delete_file("test.txt")
        self.assertIs(os.path.exists("test.txt"), False)

    def test_delete_file_doesnt_exist(self):
        delete_file("test.txt")

    def test_write_json_file(self):
        write_json(data=self.people, output_file="test.json")
        self.assertIs(os.path.exists("test.json"), True)

    def test_write_csv_file(self):
        write_csv(data=self.errors, output_file="test.csv")
        self.assertIs(os.path.exists("test.csv"), True)

    def tearDown(self) -> None:
        delete_file("test.json")
        delete_file("test.csv")
