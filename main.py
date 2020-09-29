import logging

from src.cli import CLI, parse_args

logging.getLogger().setLevel(logging.INFO)  # TODO: add parameter to toggle logging level

if __name__ == '__main__':
    # Setup CLI arguments
    config = parse_args()
    # Instantiate a CLI Class
    cli = CLI(input_directory=config.input_directory, output_directory=config.output_directory,
              error_directory=config.error_directory)

    # Run the CLI program: Control + C to quit
    cli.run()
