package main

import (
	"flag"
	"fmt"
)

func parseArgs(
	inDir, outDir, errDir, logFile *string,
	verbose, distributed, careful *bool,
	maxConcurrentFiles, maxQueuedRecords *int,
) bool {
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
	flag.IntVar(
		maxConcurrentFiles,
		"mf",
		100,
		fmt.Sprintf(
			"%s %s",
			"Controls the size of the file processing queue the program to maintain.",
			"Should correspond with the maximum amount of files being processed at once.",
		),
	)
	flag.IntVar(
		maxQueuedRecords,
		"mr",
		1000,
		fmt.Sprintf(
			"%s %s",
			"Controls the amount of queued records each channel is allowed to main.",
			"Should correspond with the maximum amount of records per file.",
		),
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
	var maxFiles, maxRecords int
	if parseArgs(
		&inDir,
		&outDir,
		&errDir,
		&logFile,
		&verbose,
		&distributed,
		&careful,
		&maxFiles,
		&maxRecords,
	) {
		return
	}

	// NOTE: will need a dictionary keeping track of the size (num lines / records) in each file
	// so that when I do things out of sequence I know when a file has been fully processed. This
	// dictionary will need locks around it because multiple processes will want to edit it.

	// Use the dictionary above with a thread-safe progressbar that reports progress for each
	// file!
	// Can give the dictionary additional fields with more information if want to track progress
	// more granually.

	// NOTE: will want to keep a "queue" of files to put in the pipeline between DirMonitor and the
	// rest of the system. The size of this queue determines the maximum number of files a user
	// will want to introduce to the program (via either starting in a directory with files or
	// copying files into the input directory) at a time, not exceeding this number.

	// NOTE: will need to convert between ResultContext and Export objects

	// NOTE: the queue mentioned above can also be used to induce the reprocessing of files once
	// completed as it will have to be aware of processes completing to limit system load.

	// Launch the goroutines
	// Wait for user input

	// Kill the goroutines -- nicely (hard kill Ctrl-C)
	// NOTE: will need to keep track of the number of goroutines listening to the killChan, then
	// send the correct number of kill signals. This will also have the effect of block the main
	// program from exiting until all the subprocesses have been killed.

	// NOTE: want to make sure that there is a clean path for any "fatal" errors. If a fatal error
	// is detected begin "nice" shutdown.
	// Exit
}
