package main

import (
	"errors"
	"sync"
)

// ProcessingContext contains the information necessary to track the processing of a file.
type ProcessingContext struct {
	NumRecords, NumProcessedRecords int
	reprocess                       bool
}

// LockedProcessingContextMap contains a lock along with a map of filenames to ProcessingContext.
type LockedProcessingContextMap struct {
	lock sync.Mutex
	Map  map[string]ProcessingContext
}

// AddProcessingFile inserts a new map entry for the provided filename with default values.
func (t *LockedProcessingContextMap) AddProcessingFile(name string) error {
	return errors.New("unimplemented")
}

// CurrentlyProcessing returns true if the provided filename is already in the Map, marking
// the file in question for reprocessing. If the filename is not in the Map, it returns false.
func (t *LockedProcessingContextMap) CurrentlyProcessing(name string) (bool, error) {
	return false, errors.New("unimplemented")
}
