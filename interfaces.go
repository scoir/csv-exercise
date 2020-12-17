package main

// ProcessingContextInterface provides a mock target for LockedProcessingContextMap
type ProcessingContextInterface interface {
	AddProcessingFile(name string) error
	CurrentlyProcessing(name string) (bool, error)
}
