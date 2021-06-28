package main

import (
	"sync"
	"testing"
)

func setup() LockedProcessingContextMap {
	process := make(chan string, 5)
	delete := make(chan string, 5)
	fileQueue := make(chan string, 5)
	contextMap := make(map[string]ProcessingContext)
	return LockedProcessingContextMap{
		maxConcurrent: 1,
		lock:          sync.Mutex{},
		process:       process,
		delete:        delete,
		fileQueue:     fileQueue,
		Map:           contextMap,
	}
}

func TestAddProcessingFile(t *testing.T) {
	obj := setup()
	obj.fileQueue <- "inqueue"
	obj.AddProcessingFile("test")
	select {
	case val := <-obj.fileQueue:
		if val != "test" {
			t.Error("Expected: test Got:", val)
		}
	}

	if !obj.CurrentlyProcessing("inqueue") {
		t.Error("inqueue not moved to processing")
	}
}
