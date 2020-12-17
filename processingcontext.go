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
	lock              sync.Mutex
	reprocess, delete chan string
	Map               map[string]ProcessingContext
}

// ProcessingContextInterface provides a mock target for LockedProcessingContextMap
type ProcessingContextInterface interface {
	AddProcessingFile(name string) error
	CurrentlyProcessing(name string) (bool, error)
	ProcessingComplete(name string) (bool, error)
	IncrementRecordsProcessed(name string) error
	SetNumRecords(name string, num int) error
}

// AddProcessingFile inserts a new map entry for the provided filename with default values.
// NOTE: This method can also be used to control the number of concurrently processing files, which
// due to current design deficiencies in etl/JSONExporter should not exceed the number of available
// processes.
func (t *LockedProcessingContextMap) AddProcessingFile(name string) error {
	return errors.New("unimplemented")
}

// CurrentlyProcessing returns true if the provided filename is already in the Map, marking
// the file in question for reprocessing. If the filename is not in the Map, it returns false.
func (t *LockedProcessingContextMap) CurrentlyProcessing(name string) (bool, error) {
	return false, errors.New("unimplemented")
}

// ProcessingComplete returns true if the provided filename has an equivalent number of records and
// processed records. Otherwise, it returns false. If it returns true, send the filename to either
// reprocess or delete, depending on the prescence of the reprocess flag. Then it removes that
// entry from the map.
func (t *LockedProcessingContextMap) ProcessingComplete(name string) (bool, error) {
	return false, errors.New("unimplemented")
}
