package main

import (
	"flag"
	"fmt"
)

func parseArgs(inDir, outDir, errDir, logFile *string, verbose, distributed, careful *bool) bool {
	flag.StringVar(inDir, "i", "", "The input directory")
	flag.StringVar(outDir, "o", "", "The output directory")
	flag.StringVar(errDir, "e", "", "The error directory")
	flag.StringVar(logFile, "l", "", "The file logs are written to")
	flag.BoolVar(verbose, "v", false, "Increases reporting when present")
	flag.BoolVar(distributed, "d", false, "Runs the program in multi-threaded mode")
	flag.BoolVar(
		careful,
		"c",
		false,
		"Runs the program in 'careful' mode, requesting additional user input",
	)

	help := flag.Bool("h", false, "Displays a usage message")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
	}

	return *help
}

func main() {
	fmt.Println("Hello World!")
	defer fmt.Println("Exiting...")

	var inDir, outDir, errDir, logFile string
	var verbose, distributed, careful bool
	if parseArgs(&inDir, &outDir, &errDir, &logFile, &verbose, &distributed, &careful) {
		return
	}

	// NOTE: will need a dictionary keeping track of the size (num lines / records) in each file
	// so that when I do things out of sequence I know when a file has been fully processed. This
	// dictionary will need locks around it because multiple processes will want to edit it.

	// Use the dictionary above with a thread-safe progressbar that reports progress for each
	// file!
	// Can give the dictionary additional fields with more information if want to track progress
	// more granually.

	// Launch the goroutines
	// Wait for user input

	// Kill the goroutines -- nicely (hard kill Ctrl-C)
	// Exit
}
