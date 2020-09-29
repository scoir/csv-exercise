import json
import logging


def write_json(data, output_file):
    """
    :param data: List of Objects that implements a toJSON method
    :param output_file: path of file to write out to
    :return:
    """
    logging.info(f"Writing json to: {output_file}")
    json_string = json.dumps([obj.toJSON() for obj in data], indent=4)
    with open(output_file, 'w') as output:
        output.write(json_string)