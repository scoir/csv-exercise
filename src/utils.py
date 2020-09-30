import csv
import json
import logging
import os


def write_json(data, output_file):
    """
    :param data: List of Objects that implements a to_json method
    :param output_file: path of file to write out to
    :return:
    """
    logging.info(f"Writing json to: {output_file}")
    json_string = json.dumps([obj.to_json() for obj in data], indent=4)
    with open(output_file, 'w') as output:
        output.write(json_string)


def delete_file(target_file: str):
    """
    helper to attempt to delete a file.
    :param target_file: location of file to delete
    :return:
    """
    logging.info(f"Deleting {target_file}")
    if os.path.exists(target_file):
        os.remove(target_file)


def create_directory_if_not_exists(directory: str):
    """
    :param directory: directory to create if it does not already exist
    :return:
    """
    logging.info(f"Checking to see if {directory} exists")
    if not os.path.exists(directory):
        logging.info(f"Creating directory: {directory}")
        os.mkdir(directory)


def write_csv(data, output_file):
    """
    write_csv will take a given list and write it out to an output file. Each row object needs to implement a `to_csv`
    method.
    :param data: list of objects with a to_csv method
    :param output_file: location of file to write data to
    :return:
    """
    with open(output_file, mode="w") as csv_file:
        writer = csv.writer(csv_file, delimiter=",")
        for row in data:
            # TODO: Not particularly happy about this being hardcoded. This logic should live in the Error class
            writer.writerow(row.to_csv())
