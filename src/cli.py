import csv
import logging

from src.models import Person, Error
from src.utils import write_json


def run(input_directory, output_directory, error_directory):
    """
    run will listen for changes on the input output_directory and process events
    when new files are created.
    :param input_directory:
    :param output_directory:
    :param error_directory:
    :return:
    """

    logging.info(f"Input: {input_directory}, output: {output_directory}, errors: {error_directory}")
    # instantiate watchdog observer: https://pypi.org/project/watchdog/
    # observer = Observer()
    # observer.schedule(CsvFileWatcher, input_directory, recursive=True)
    # observer.start()
    # try:
    #     while True:
    #         time.sleep(1)
    # except KeyboardInterrupt:
    #     observer.stop()
    #
    # observer.join()
    people, errors = process_new_input_file(f"{input_directory}/MOCK_DATA.csv")  # TODO: Refactor to observer pattern
    write_json(people, f"{output_directory}/MOCK_DATA.json") # TODO: Refactor to get filename from change event
    write_json(errors, f"{error_directory}/MOCK_DATA.json") # TODO: Refactor to get filename from change event


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
                person = Person(id=row['INTERNAL_ID'], first_name=row['FIRST_NAME'], last_name=row['LAST_NAME'],
                                middle_name=row['MIDDLE_NAME'], phone_number=row['PHONE_NUM'])
                people.append(person)
            except ValueError as e:
                logging.error(f"Error processing row: {e}")
                error = Error(row=row_number, error=e)
                errors.append(error)

    return people, errors