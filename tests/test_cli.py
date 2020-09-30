import os
import unittest
from unittest import TestCase

from src.cli import CLI


class TestCLI(TestCase):
    def setUp(self):
        self.cli = CLI(input_directory="test_input", output_directory="test_output", error_directory="test_error")

    def test_cli_initialized_with_defaults(self):
        self.assertIsInstance(self.cli, CLI)

    def tearDown(self) -> None:
        test_directories = [self.cli.input_directory, self.cli.output_directory, self.cli.error_directory]
        for directory in test_directories:
            os.rmdir(directory)


if __name__ == '__main__':
    unittest.main()