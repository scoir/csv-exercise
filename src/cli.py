import argparse
import csv
import logging
import os
import time

from watchdog.events import PatternMatchingEventHandler
from watchdog.observers import Observer

from src.models import Person, Error
from src.utils import write_json, delete_file, create_directory_if_not_exists, write_csv


class CLI:
    """
    CLI represents a Command Line Interface object. CLI is responsible for
    creating directories, tracking processed files, and handling on_create events.
    """

    def __init__(self, input_directory="input_directory", output_directory="output_directory",
                 error_directory="error_directory"):
        self.input_directory = input_directory
        self.output_directory = output_directory
        self.error_directory = error_directory
        self.event_handler = PatternMatchingEventHandler(patterns="*.csv", )
        self.event_handler.on_created = self._on_created
        self.processed = set()
        self._create_directories_if_not_exists()

    def _create_directories_if_not_exists(self):
        """
        _create_directories_if_not_exists is an internal helper method that will ensure each
        directory exists prior to executing run.
        :return:
        """
        directories = [self.input_directory, self.output_directory, self.error_directory]
        for directory in directories:
            create_directory_if_not_exists(directory)

    def run(self):
        """
        run will listen for changes on the input output_directory and process events
        when new files are created.
        :return:
        """

        logging.info(f"Input: {self.input_directory}, output: {self.output_directory}, errors: {self.error_directory}")
        observer = Observer()
        observer.schedule(self.event_handler, self.input_directory, recursive=True)
        observer.start()
        logging.info(f"Preparing to observe {self.input_directory} for .csv files. Control + C to exit")
        try:
            while True:
                time.sleep(1)
        except KeyboardInterrupt:
            # TODO: Print out each of the files processed
            observer.stop()
            observer.join()

    def _on_created(self, event):
        """
        _on_created is a Watchdog Event Handler. It will take the filename that
        was created, process the input file, and write the result out to
        the respective files.
        :param event:
        :return:
        """
        logging.info(f"handling on_create event: {event}")
        src_path = event.src_path
        if src_path in self.processed:
            logging.info(f"File {src_path} has already been processed.")
            self._delete_file(target_path=src_path)
            # TODO: Is this an appropriate way of breaking from the method?
            return
        else:
            self.processed.add(src_path)
        people, errors = self.process_new_input_file(
                f"{src_path}")

        self._delete_file(target_path=src_path)
        _, file_name = os.path.split(src_path)
        target_path = file_name.replace(".csv", ".json")
        people_file = f"{self.output_directory}/{target_path}"
        errors_file = f"{self.error_directory}/{file_name}"

        self._write_json_to_file(data=people, output_file=people_file)
        self._write_csv_to_file(data=errors, output_file=errors_file)

        logging.info(f"Successfully parsed out {src_path}!: {event}")

    @staticmethod
    def _write_csv_to_file(data, output_file):
        logging.info(f"Writing out to : {output_file}")
        write_csv(data=data, field_names=['LINE_NUM', 'ERROR_MSG'], output_file=output_file)

    @staticmethod
    def _write_json_to_file(data, output_file):
        """
        _write_json_to_file is a helper function to call write_json
        :param data: JSON Serializable object to write out
        :param output_file: path of the file to write data to
        :return:
        """
        logging.info(f"Writing out to: {output_file}")
        write_json(data, output_file=output_file)

    @staticmethod
    def _delete_file(target_path):
        """
        _delete_file is a private helper function that calls delete file.
        :param target_path: Path of the file to delete
        :return:
        """
        delete_file(target_path)

    @staticmethod
    def process_new_input_file(file: str) -> tuple:
        """
        :rtype: tuple
        :param file:
        :return:
        """
        people = []
        errors = []
        logging.info(f"Received request to process file: {file}")
        with open(file, newline='') as csv_file:
            reader = csv.DictReader(csv_file)
            for row_number, row in enumerate(reader):
                try:
                    person = Person(internal_id=row['INTERNAL_ID'], first_name=row['FIRST_NAME'],
                                    last_name=row['LAST_NAME'],
                                    middle_name=row['MIDDLE_NAME'], phone_number=row['PHONE_NUM'])
                    people.append(person)
                except ValueError as e:
                    logging.error(f"Error processing row: {e}")
                    error = Error(row=row_number, error=e)
                    errors.append(error)

        return people, errors


def parse_args() -> {}:
    """
    Scaffold the CLI tool and parse_args defaults.
    :return: Object with configuration values
    """
    parser = argparse.ArgumentParser(description=DEFAULTS['DESCRIPTION'])
    parser.add_argument('-i', '--input-directory', help="Input Directory to watch", default=DEFAULTS['INPUT_DIRECTORY'])
    parser.add_argument('-o', '--output-directory', help="Directory to write our results",
                        default=DEFAULTS['OUTPUT_DIRECTORY'])
    parser.add_argument('-e', '--error-directory', help="Directory to write out errors",
                        default=DEFAULTS['ERROR_DIRECTORY'])
    args = parser.parse_args()
    return args


DEFAULTS = {
    'INPUT_DIRECTORY' : 'input_directory',
    'OUTPUT_DIRECTORY': 'output_directory',
    'ERROR_DIRECTORY' : 'error_directory',
    'DESCRIPTION'     : 'Scoir CSV Exercise. Watches a specified input directory and acts upon CSV Files being '
                        'created. '
                        'Control-C to quit execution',
}
