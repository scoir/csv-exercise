package main

import (
	"github.com/8bitstout/csv-exercise/parser"
	"github.com/urfave/cli"
	"log"
	"os"
)

const (
	DEFAULT_ROW_LENGTH       = 5
	DEFAULT_INPUT_DIRECTORY  = "./input"
	DEFAULT_OUTPUT_DIRECTORY = "./output"
	DEFAULT_ERRORS_DIRECTORY = "./errors"
)

func main() {
	app := &cli.App{
		Name:  "csv-exercise",
		Usage: "Parse a csv to json",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "input-dir",
				Value: DEFAULT_INPUT_DIRECTORY,
				Usage: "The directory to watch for new files to be parsed",
			},
			&cli.StringFlag{
				Name:  "output-dir",
				Value: DEFAULT_OUTPUT_DIRECTORY,
				Usage: "The directory parsed files will be created in",
			},
			&cli.StringFlag{
				Name:  "error-dir",
				Value: DEFAULT_ERRORS_DIRECTORY,
				Usage: "The directory error files will be created in",
			},
		},
		Action: func(c *cli.Context) error {
			p := parser.NewParser(DEFAULT_ROW_LENGTH, DEFAULT_INPUT_DIRECTORY, DEFAULT_OUTPUT_DIRECTORY, DEFAULT_ERRORS_DIRECTORY)
			for {
				p.Watch()
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
