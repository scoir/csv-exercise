package main

import "errors"

// Logger is a daemon that writes all log messages sent to its channel to the log file.
// NOTE: the channel might block and cause hang-ups if the logger is under heavy load.
func Logger(filename *string, logChan chan string) error {
	// select statement is the core here
	// on recieving "QUIT" it shuts down
	return errors.New("not implemented")
}
