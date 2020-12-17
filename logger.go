package main

import (
	"errors"
	"fmt"
	"time"
)

// Logger is a daemon that writes all log messages sent to its channel to the log file.
// NOTE: the channel might block and cause hang-ups if the logger is under heavy load.
// NOTE: this is just a barebones approach, an object-based approach would probably be better.
func Logger(filename *string, logChan chan string) error {
	// select statement is the core here
	// on recieving "QUIT" it shuts down
	return errors.New("not implemented")
}

// LogError takes an error logs it, formatted to include the time of the message
func LogError(err error, logChan chan string) {
	t := time.Now().UTC()
	logChan <- fmt.Sprintf("%s;%s", t.Format("2006-01-02 15:04:05"), err.Error())
}
