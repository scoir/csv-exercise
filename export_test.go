package main

import (
	"os"
	"testing"
	"time"

	mock "github.com/JackStillwell/csv-exercise/mock"
	gomock "github.com/golang/mock/gomock"
)

// NOTE: the ExportManager and Exporter deserve a more comprehensive testing suite. Time
// constraints mean I need to work fast, so I'm limiting myself to a few "step-through" tests to
// give at least a small amount of coverage.

func TestExporter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	assignedFile := "test/test.in"
	input := make(chan ExportInterface, 0)
	errChan := make(chan error, 0)
	kill := make(chan bool, 0)
	m := mock.NewMockProcessingContextInterface(ctrl)
	gomock.InOrder(
		m.EXPECT().IncrementRecordsProcessed("test/test.in").Return(nil),
		m.EXPECT().ProcessingComplete("test/test.in").Return(false, nil),
		m.EXPECT().IncrementRecordsProcessed("test/test.in").Return(nil),
		m.EXPECT().ProcessingComplete("test/test.in").Return(true, nil),
	)

	go Exporter(&assignedFile, input, m, kill, errChan)

	exports := []Export{
		{
			"test/test.in",
			"test",
			&ErrorFormat{
				1,
				"testerror",
			},
		},
		{
			"test/test.in",
			"test",
			&PersonFormat{
				Person{
					1,
					PersonName{
						"a",
						"b",
						"c",
					},
					"#",
				},
			},
		},
	}

	for _, val := range exports {
		input <- &val
	}

	kill <- true

	os.Remove("test/test.csv")
	os.Remove("test/test.csv")
}

func TestExportManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	input := make(chan ExportInterface, 0)
	logger := make(chan string, 10)
	kill := make(chan bool, 0)
	m := mock.NewMockProcessingContextInterface(ctrl)
	m.EXPECT().IncrementRecordsProcessed("test/test.in").AnyTimes().Return(nil)
	m.EXPECT().IncrementRecordsProcessed("test/test2.in").AnyTimes().Return(nil)
	m.EXPECT().ProcessingComplete("test/test.in").Return(false, nil).AnyTimes()
	m.EXPECT().ProcessingComplete("test/test2.in").Return(false, nil).AnyTimes()

	go ExportManager(input, m, 2, 5, logger, kill)

	exports := []Export{
		{
			"test/test.in",
			"test",
			&ErrorFormat{
				1,
				"testerror",
			},
		},
		{
			"test/test2.in",
			"test",
			&PersonFormat{
				Person{
					1,
					PersonName{
						"a",
						"b",
						"c",
					},
					"#",
				},
			},
		},
		{
			"test/test2.in",
			"test",
			&ErrorFormat{
				1,
				"testerror",
			},
		},
		{
			"test/test.in",
			"test",
			&PersonFormat{
				Person{
					1,
					PersonName{
						"a",
						"b",
						"c",
					},
					"#",
				},
			},
		},
	}

	for idx := range exports {
		input <- &exports[idx]
	}

	time.Sleep(time.Second * 1)
	kill <- true
	time.Sleep(time.Second * 1)

	os.Remove("test/test.csv")
	os.Remove("test/test.json")
	os.Remove("test/test2.csv")
	os.Remove("test/test2.json")
}
