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


def delete_file(target_file):
    logging.info(f"Deleting {target_file}")
    os.remove(target_file)


def create_directory_if_not_exists(directory: str):
    """
    create_directory_if_not_exists will
    :param directory:
    :return:
    """
    logging.info(f"Checking to see if {directory} exists")
    if not os.path.exists(directory):
        logging.info(f"Creating directory: {directory}")
        os.mkdir(directory)


def write_csv(data, field_names, output_file):
    """
    write_csv will take a given list and write it out to an output file.
    :param data:
    :param field_names:
    :param output_file:
    :return:
    """
    with open(output_file, mode="w") as csv_file:
        writer = csv.DictWriter(csv_file, field_names)
        writer.writeheader()
        for row in data:
            # TODO: Not particularly happy about this being hardcoded. This logic should live in the Error class
            writer.writerow({"LINE_NUM": row.row, "ERROR_MSG": row.message})
