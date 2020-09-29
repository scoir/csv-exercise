import logging

# TODO: http://thepythoncorner.com/dev/how-to-create-a-watchdog-in-python-to-look-for-filesystem-changes/
from src.cli import run

DEFAULTS = {
    'INPUT_DIRECTORY' : 'data',
    'OUTPUT_DIRECTORY': 'output_directory',
    'ERROR_DIRECTORY' : 'error_directory'
}

logging.getLogger().setLevel(logging.INFO) # TODO: add parameter to turn this on


def setup() -> {}:
    """
    Scaffold the CLI tool and setup defaults.
    :return: Object with configuration values
    """
    # TODO: Add arguments for input, output and error directories
    return DEFAULTS


if __name__ == '__main__':
    # Setup CLI arguments
    config = setup()
    run(input_directory=config['INPUT_DIRECTORY'], output_directory=config['OUTPUT_DIRECTORY'],
        error_directory=config['ERROR_DIRECTORY'])
