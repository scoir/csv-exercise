package main

import (
	"errors"
)

// ExportManager sorts incoming objects and assigns them to the proper new or
// existing Exporter instance.
// NOTE: This process will be a bottleneck as it has to handle every record. Its assigned task is
// very light due to this, merely sending files to the right place.
func ExportManager(
	exportStream chan ExportInterface,
	processingContext ProcessingContextInterface,
	numProcesses, maxRecords int,
	logger chan string,
	killChan chan bool,
) {
	// Initiate exporters
	subProcessKill := make(chan bool, 0)
	subProcessError := make(chan error, numProcesses)
	// NOTE: the filenames are the source filename used in processingContext
	assignedFiles := make([]string, numProcesses)
	exportChans := make([]chan ExportInterface, numProcesses)
	for i := 0; i < numProcesses; i++ {
		exportChans[i] = make(chan ExportInterface, maxRecords)
	}

	for idx := range assignedFiles {
		go Exporter(
			&assignedFiles[idx],
			exportChans[idx],
			processingContext,
			subProcessKill,
			subProcessError,
		)
	}

	// Act on every recieved ExportInterface
	kill := false
	for !kill {
		select {
		case kill = <-killChan:
			// ensure all sub-processes are killed before exiting
			for i := 0; i < numProcesses; i++ {
				subProcessKill <- true
			}
			continue
		case export := <-exportStream:
			// fmt.Println("Manager Recieved", export.GetSource(), export.GetFileName())
			sourceFile := export.GetSource()
			// First, check to see if any running process has the file assigned
			handled := false
			for idx, val := range assignedFiles {
				if val == sourceFile {
					// fmt.Println("Manager Assigned", export.GetFileName(), "to", idx)
					exportChans[idx] <- export
					handled = true
					break
				}
			}

			// If not, assign an idle process
			if !handled {
				for idx, val := range assignedFiles {
					if val == "" {
						// fmt.Println("Manager Assigned", idx, "to", export.GetSource())
						// fmt.Println("Manager Assigned", export.GetFileName(), "to", idx)
						assignedFiles[idx] = sourceFile
						exportChans[idx] <- export
						handled = true
						break
					}
				}
			}

			// If not handled at this point, something is broken
			if !handled {
				LogError(errors.New(
					"ExportManager unable to distribute work",
				), logger)
			}

		case subErr := <-subProcessError:
			LogError(subErr, logger)
		}
	}
}
