package main

import (
	"errors"
	"io/ioutil"
	"path"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// DirMonitor keeps track of the input directory and flags files for processing.
func DirMonitor(
	inDir *string,
	processingMap ProcessingContextInterface,
	toProcess,
	logger chan string,
	killChan chan bool,
	errChan chan error,
) {

	// First, get lists of all the files in each directory
	files, err := ioutil.ReadDir(*inDir)
	if err != nil {
		errChan <- err
		return
	}

	// If the files are CSV files, process them
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".csv") {
			continue
		}
		fileName := path.Join(*inDir, file.Name())
		processFile(fileName, processingMap, toProcess, logger)
	}

	// Set up a Watcher for any new files
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		errChan <- err
		return
	}
	defer watcher.Close()

	err = watcher.Add(*inDir)
	if err != nil {
		errChan <- err
		return
	}

	kill := false
	for !kill {
		select {
		case kill = <-killChan:
			continue
		case event, ok := <-watcher.Events:
			if !ok {
				errChan <- errors.New("Fatal OK Error in fsnotify watcher Events")
				return
			}

			if fileModified(event.Op) && strings.HasSuffix(event.Name, ".csv") {
				fileName := event.Name

				if processing, err := processingMap.CurrentlyProcessing(fileName); processing || err != nil {
					if err != nil {
						LogError(err, logger)
					}
					continue
				}
				processFile(fileName, processingMap, toProcess, logger)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				errChan <- errors.New("Fatal OK Error in fsnotify watcher Errors")
				return
			}
			errChan <- err
			return
		}
	}
}

func processFile(fileName string, processingMap ProcessingContextInterface, toProcess, logger chan string) {
	pMErr := processingMap.AddProcessingFile(fileName)
	if pMErr != nil {
		LogError(pMErr, logger)
		return
	}
	toProcess <- fileName
}

func fileModified(op fsnotify.Op) bool {
	return op == fsnotify.Create || op == fsnotify.Write || op == fsnotify.Rename
}
