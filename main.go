package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
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

	if inDir == "" || outDir == "" || errDir == "" || logFile == "" {
		fmt.Println("The arg flags i o e and l are required")
		return
	}

	// toProcess := make(chan string, maxRecords)
	kill := make(chan bool)
	logger := make(chan string, maxRecords)
	errChan := make(chan error, maxRecords)
	processFile := make(chan string, maxFiles)
	deleteFile := make(chan string, maxFiles)
	fileQueue := make(chan string, maxFiles)
	contextMap := make(map[string]ProcessingContext)
	processingContext := LockedProcessingContextMap{
		maxConcurrent: maxFiles,
		lock:          sync.Mutex{},
		process:       processFile,
		fileQueue:     fileQueue,
		delete:        deleteFile,
		Map:           contextMap,
	}

	// Get number of processes available to the runtime
	numCPUs := runtime.NumCPU()

	go Logger(&logFile, logger)

	go DirMonitor(&inDir, &processingContext, logger, kill, errChan)

	rawData := make(chan ResultContext, maxRecords)

	resultData := make(chan ResultContext)
	verifyChan := make(chan PersonContext)

	// split between errors and non-errors here
	go errSplitter(rawData, resultData, verifyChan, kill)

	for i := 0; i < numCPUs; i++ {
		go CSVImporter(processFile, rawData, &processingContext, logger, kill)
		go Verifier(verifyChan, resultData, kill)
	}

	// convert between resultcontext and Export here
	exportData := make(chan ExportInterface)
	go resultToExport(&outDir, &errDir, resultData, exportData, kill)
	go ExportManager(exportData, &processingContext, numCPUs, maxRecords, logger, kill)

	go fileDeleter(deleteFile, kill)

	// now waiting for user input to kill
	var input string
	fmt.Println("Hit the Enter Key to Kill the Program (or Ctrl-C for hard kill)")
	fmt.Scanln(&input)

	// after recieve user input, issue the kill signals and die once they've all been recieved
	numKillSignals := (2 * numCPUs) + 5 // all importers, verifiers, dirmonitor, exportmanager,
	//  converter, splitter, deleter, and logger
	for i := 0; i < numKillSignals; i++ {
		kill <- true
	}

	time.Sleep(time.Second)

	fmt.Println("Exited Successfully")
	return
}

func errSplitter(raw, results chan ResultContext, verifyChan chan PersonContext, killChan chan bool) {
	kill := false
	for !kill {
		select {
		case kill = <-killChan:
			continue
		case rawData := <-raw:
			if rawData.result == "" {
				verifyChan <- rawData.context
			} else {
				results <- rawData
			}
		}
	}
}

func resultToExport(
	outDir, errDir *string,
	results chan ResultContext,
	exports chan ExportInterface,
	killChan chan bool,
) {
	kill := false
	for !kill {
		select {
		case kill = <-killChan:
			continue
		case r := <-results:
			if r.result == "" {
				exports <- &Export{
					source:  r.context.FileName,
					destDir: *outDir,
					obj:     &PersonFormat{r.context.Record},
				}
			} else {
				exports <- &Export{
					source:  r.context.FileName,
					destDir: *errDir,
					obj:     &ErrorFormat{r.context.LineNumber, r.result},
				}
			}
		}
	}
}

func fileDeleter(fileChan chan string, killChan chan bool) {
	kill := false
	for !kill {
		select {
		case kill = <-killChan:
			continue
		case f := <-fileChan:
			os.Remove(f)
		}
	}
}
