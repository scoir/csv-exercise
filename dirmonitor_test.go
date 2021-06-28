package main

import (
	"os"
	"testing"
	"time"

	mock "github.com/JackStillwell/csv-exercise/mock"
	gomock "github.com/golang/mock/gomock"
)

func TestFindsNewFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockProcessingContextInterface(ctrl)
	m.EXPECT().CurrentlyProcessing("test/testnew.csv").Return(false)
	m.EXPECT().AddProcessingFile("test/testnew.csv")
	logger := make(chan string, 1)
	kill := make(chan bool, 1)
	errChan := make(chan error, 1)
	testDir := "test"
	go DirMonitor(&testDir, m, logger, kill, errChan)
	time.Sleep(time.Second * 1)
	if cErr := createFile("test/testnew.csv"); cErr != nil {
		t.Error(cErr)
	}

	time.Sleep(time.Second * 1)

	if rErr := os.Remove("test/testnew.csv"); rErr != nil {
		t.Error(rErr)
	}

	kill <- true

	// the test is in the mock
}

func TestFindsExistingFile(t *testing.T) {
	if cErr := createFile("test/test.csv"); cErr != nil {
		t.Error(cErr)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockProcessingContextInterface(ctrl)
	m.EXPECT().AddProcessingFile("test/test.csv")
	logger := make(chan string, 1)
	errChan := make(chan error, 1)
	kill := make(chan bool, 1)
	testDir := "test"
	go DirMonitor(&testDir, m, logger, kill, errChan)

	time.Sleep(time.Second)

	if rErr := os.Remove("test/test.csv"); rErr != nil {
		t.Error(rErr)
	}

	kill <- true

	// the test is in the mock
}

func TestIgnoresNonCSVFiles(t *testing.T) {
	if cErr := createFile("test/test.txt"); cErr != nil {
		t.Error(cErr)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockProcessingContextInterface(ctrl)
	logger := make(chan string, 1)
	errChan := make(chan error, 1)
	kill := make(chan bool, 1)
	testDir := "test"
	go DirMonitor(&testDir, m, logger, kill, errChan)
	time.Sleep(time.Second * 1)
	if cErr := createFile("test/testnew.txt"); cErr != nil {
		t.Error(cErr)
	}

	time.Sleep(time.Second)

	if rErr := os.Remove("test/test.txt"); rErr != nil {
		t.Error(rErr)
	}

	if rErr := os.Remove("test/testnew.txt"); rErr != nil {
		t.Error(rErr)
	}

	kill <- true

	// the test is in the mock
}

func TestMarksModifiedFile(t *testing.T) {
	if cErr := createFile("test/test.csv"); cErr != nil {
		t.Error(cErr)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockProcessingContextInterface(ctrl)
	m.EXPECT().CurrentlyProcessing("test/test.csv").Return(true)
	m.EXPECT().AddProcessingFile("test/test.csv")
	logger := make(chan string, 1)
	errChan := make(chan error, 1)
	kill := make(chan bool, 1)
	testDir := "test"
	go DirMonitor(&testDir, m, logger, kill, errChan)
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

func createFile(name string) error {
	f, createErr := os.Create(name)
	if createErr != nil {
		return createErr
	}
	f.Close()
	return nil
}
