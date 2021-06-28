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
	maxConcurrent              int
	lock                       sync.Mutex
	process, delete, fileQueue chan string
	Map                        map[string]ProcessingContext
}

// ProcessingContextInterface provides a mock target for LockedProcessingContextMap
type ProcessingContextInterface interface {
	AddProcessingFile(name string)
	CurrentlyProcessing(name string) bool
	ProcessingComplete(name string) (bool, error)
	IncrementRecordsProcessed(name string) error
	SetNumRecords(name string, num int) error
	NumCurrentlyProcessing() int
}

// AddProcessingFile inserts a new map entry for the provided filename with default values.
// NOTE: This method can also be used to control the number of concurrently processing files, which
// should not exceed the number of available processes.
func (t *LockedProcessingContextMap) AddProcessingFile(name string) {
	t.lock.Lock()
	t.fileQueue <- name
	t.updateQueue()
	t.lock.Unlock()
}

// CurrentlyProcessing returns true if the provided filename is already in the Map, marking the
// file in question for reprocessing. If the filename is not in the Map, it returns false.
func (t *LockedProcessingContextMap) CurrentlyProcessing(name string) bool {
	t.lock.Lock()
	if v, ok := t.Map[name]; ok {
		v.reprocess = true
		t.lock.Unlock()
		return true
	}
	t.lock.Unlock()
	return false
}

// ProcessingComplete returns true if the provided filename has an equivalent number of records and
// processed records. Otherwise, it returns false. If it returns true, send the filename to either
// reprocess or delete, depending on the prescence of the reprocess flag. Then it removes that
// entry from the map.
func (t *LockedProcessingContextMap) ProcessingComplete(name string) (bool, error) {
	t.lock.Lock()
	if v, ok := t.Map[name]; ok {
		if v.NumRecords == v.NumProcessedRecords {
			if v.reprocess {
				t.fileQueue <- name
			} else {
				t.delete <- name
			}
			delete(t.Map, name)
			t.updateQueue()
			t.lock.Unlock()
			return true, nil
		}
		t.lock.Unlock()
		return false, nil
	}
	t.lock.Unlock()
	return false, errors.New("ProcessingComplete cannot find file " + name + " in the processing map")
}

// IncrementRecordsProcessed increases the "NumProcessedRecords" field of a ProcessingContext by 1
func (t *LockedProcessingContextMap) IncrementRecordsProcessed(name string) error {
	t.lock.Lock()
	if v, ok := t.Map[name]; ok {
		v.NumProcessedRecords++
		t.lock.Unlock()
		return nil
	}
	t.lock.Unlock()
	return errors.New("ProcessingComplete cannot find file " + name + " in the processing map")
}

// SetNumRecords sets the "NumRecords" field of a ProcessingContext
func (t *LockedProcessingContextMap) SetNumRecords(name string, num int) error {
	t.lock.Lock()
	if v, ok := t.Map[name]; ok {
		v.NumRecords = num
		t.lock.Unlock()
		return nil
	}
	t.lock.Unlock()
	return errors.New("ProcessingComplete cannot find file " + name + " in the processing map")
}

// NumCurrentlyProcessing returns the amount of items currently present in the ProcessingContextMap
func (t *LockedProcessingContextMap) NumCurrentlyProcessing() int {
	t.lock.Lock()
	numItems := len(t.Map)
	t.lock.Unlock()
	return numItems
}

// updateQueue checks the amount of currently processing files and kicks off more files if they are
// available in the fileQueue.
func (t *LockedProcessingContextMap) updateQueue() {
	queueHasItems := true
	for len(t.Map) < t.maxConcurrent && queueHasItems {
		select {
		case file := <-t.fileQueue:
			// make sure its not in the map
			if v, ok := t.Map[file]; ok {
				v.reprocess = true
				continue
			}
			// add it
			t.Map[file] = ProcessingContext{0, 0, false}

			// send it forward
			t.process <- file
		default:
			queueHasItems = false
		}
	}
}
