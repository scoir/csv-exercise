import logging

from src.cli import CLI, parse_args


if __name__ == '__main__':
    # Setup CLI arguments
    config = parse_args()
    if config.verbose:
        logging.getLogger().setLevel(logging.INFO)
    # Instantiate a CLI Class
    cli = CLI(input_directory=config.input_directory, output_directory=config.output_directory,
              error_directory=config.error_directory)

    # Run the CLI program: Control + C to quit
    cli.run()
