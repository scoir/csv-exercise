package main

import (
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	mock "github.com/JackStillwell/csv-exercise/mock"
	gomock "github.com/golang/mock/gomock"
)

func TestFindsNewFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockProcessingContextInterface(ctrl)
	m.EXPECT().CurrentlyProcessing("test/testnew.csv").Return(false, nil)
	m.EXPECT().AddProcessingFile("test/testnew.csv").Return(nil)
	toProcess := make(chan string, 1)
	logger := make(chan string, 1)
	kill := make(chan bool, 1)
	errChan := make(chan error, 1)
	testDir := "test"
	go DirMonitor(&testDir, m, toProcess, logger, kill, errChan)
	time.Sleep(time.Second * 1)
	if cErr := createFile("test/testnew.csv"); cErr != nil {
		t.Error(cErr)
	}
	var val string
	select {
	case <-time.After(time.Second * 5):
	case val = <-toProcess:
	}

	if rErr := os.Remove("test/testnew.csv"); rErr != nil {
		t.Error(rErr)
	}

	kill <- true

	if val != "test/testnew.csv" {
		t.Errorf("Expected: test/testnew.csv Got: %s", val)
	}
}

func TestFindsExistingFile(t *testing.T) {
	if cErr := createFile("test/test.csv"); cErr != nil {
		t.Error(cErr)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockProcessingContextInterface(ctrl)
	m.EXPECT().AddProcessingFile("test/test.csv").Return(nil)
	toProcess := make(chan string, 1)
	logger := make(chan string, 1)
	errChan := make(chan error, 1)
	kill := make(chan bool, 1)
	testDir := "test"
	go DirMonitor(&testDir, m, toProcess, logger, kill, errChan)
	var val string
	select {
	case <-time.After(time.Second * 5):
	case val = <-toProcess:
	}

	if rErr := os.Remove("test/test.csv"); rErr != nil {
		t.Error(rErr)
	}

	kill <- true

	if val != "test/test.csv" {
		t.Errorf("Expected: test/test.csv Got: %s", val)
	}
}

func TestIgnoresNonCSVFiles(t *testing.T) {
	if cErr := createFile("test/test.txt"); cErr != nil {
		t.Error(cErr)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockProcessingContextInterface(ctrl)
	toProcess := make(chan string, 1)
	logger := make(chan string, 1)
	errChan := make(chan error, 1)
	kill := make(chan bool, 1)
	testDir := "test"
	go DirMonitor(&testDir, m, toProcess, logger, kill, errChan)
	time.Sleep(time.Second * 1)
	if cErr := createFile("test/testnew.txt"); cErr != nil {
		t.Error(cErr)
	}
	var val string
	select {
	case <-time.After(time.Second * 1):
	case val = <-toProcess:
	}

	if rErr := os.Remove("test/test.txt"); rErr != nil {
		t.Error(rErr)
	}

	if rErr := os.Remove("test/testnew.txt"); rErr != nil {
		t.Error(rErr)
	}

	kill <- true

	if val != "" {
		t.Errorf("Expected:  Got: %s", val)
	}
}

func TestMarksModifiedFile(t *testing.T) {
	if cErr := createFile("test/test.csv"); cErr != nil {
		t.Error(cErr)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockProcessingContextInterface(ctrl)
	m.EXPECT().CurrentlyProcessing("test/test.csv").Return(true, nil)
	m.EXPECT().AddProcessingFile("test/test.csv").Return(nil)
	toProcess := make(chan string, 1)
	logger := make(chan string, 1)
	errChan := make(chan error, 1)
	kill := make(chan bool, 1)
	testDir := "test"
	go DirMonitor(&testDir, m, toProcess, logger, kill, errChan)
	time.Sleep(time.Second * 1)
	f, oErr := os.OpenFile("test/test.csv", os.O_RDWR, 0644)
	if oErr != nil {
		t.Error(oErr)
	}
	f.WriteString("test")
	f.Close()

	time.Sleep(time.Second * 1)
	if rErr := os.Remove("test/test.csv"); rErr != nil {
		t.Error(rErr)
	}

	kill <- true

	// The check here is in the mock - only 1 addprocessing and a currentlyprocessing
}

func TestLogsProcessingMapError(t *testing.T) {
	if cErr := createFile("test/test.csv"); cErr != nil {
		t.Error(cErr)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockProcessingContextInterface(ctrl)
	m.EXPECT().AddProcessingFile("test/test.csv").Return(errors.New("test"))
	toProcess := make(chan string, 1)
	logger := make(chan string, 1)
	errChan := make(chan error, 1)
	kill := make(chan bool, 1)
	testDir := "test"
	go DirMonitor(&testDir, m, toProcess, logger, kill, errChan)
	time.Sleep(time.Second * 1)
	if rErr := os.Remove("test/test.csv"); rErr != nil {
		t.Error(rErr)
	}

	var val string
	select {
	case <-time.After(time.Second * 1):
	case val = <-logger:
	}

	kill <- true

	if !strings.HasSuffix(val, "test") {
		t.Errorf("Expected: test Got: %s", val)
	}
}

func createFile(name string) error {
	f, createErr := os.Create(name)
	if createErr != nil {
		return createErr
	}
	f.Close()
	return nil
}
