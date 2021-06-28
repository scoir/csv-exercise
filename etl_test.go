package main

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	mock "github.com/JackStillwell/csv-exercise/mock"
	gomock "github.com/golang/mock/gomock"
)

// NOTE: The CSVImporter stack is too monolithic and needs to be refactored. As a consequence,
// the tests are pretty knarly.

func csvImporterSetup(dir string) error {
	fC1 := "INTERNAL_ID,FIRST_NAME,MIDDLE_NAME,LAST_NAME,PHONE_NUM\n" +
		"1,,,,\n"

	fC2 := "x\n"

	fC3 := fC1 + ",,,\n"

	fC4 := fC1 + "a,,,,\n"

	for idx, val := range [4]string{fC1, fC2, fC3, fC4} {
		filepath := path.Join(dir, fmt.Sprintf("test%d.csv", idx))
		f, err := os.Create(filepath)
		if err != nil {
			return err
		}
		f.WriteString(val)
		f.Close()
	}
	return nil
}

func csvImporterCleanup(dir string) error {
	for idx := 0; idx < 4; idx++ {
		filepath := path.Join(dir, fmt.Sprintf("test%d.csv", idx))
		err := os.Remove(filepath)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestCSVToPerson(t *testing.T) {
	csvImporterSetup("test")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := ResultContext{
		"",
		PersonContext{"test/test0.csv", 2, Person{1, PersonName{"", "", ""}, ""}},
	}
	input := make(chan string, 0)
	output := make(chan ResultContext, 0)
	kill := make(chan bool, 0)
	m := mock.NewMockProcessingContextInterface(ctrl)
	m.EXPECT().SetNumRecords("test/test0.csv", 1).Return(nil)

	go CSVImporter(input, output, m, input, kill)
	input <- "test/test0.csv"

	retVal := ResultContext{}
	select {
	case retVal = <-output:
	case logString := <-input:
		kill <- true
		t.Error("Received Error:", logString)
	case <-time.After(time.Second * 5):
	}

	kill <- true

	if retVal != expected {
		t.Error("Expected:", expected, "Got:", retVal)
	}

	csvImporterCleanup("test")
}

func TestCSVToPersonBadHeader(t *testing.T) {
	csvImporterSetup("test")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	header := "INTERNAL_ID,FIRST_NAME,MIDDLE_NAME,LAST_NAME,PHONE_NUM"
	expected := ResultContext{
		"Incorrect CSV Header | Expected: " + header + " | Got: x",
		PersonContext{"test/test1.csv", 1, Person{0, PersonName{"", "", ""}, ""}},
	}
	input := make(chan string, 0)
	output := make(chan ResultContext, 0)
	kill := make(chan bool, 0)
	m := mock.NewMockProcessingContextInterface(ctrl)

	go CSVImporter(input, output, m, input, kill)
	input <- "test/test1.csv"

	retVal := ResultContext{}
	select {
	case retVal = <-output:
	case logString := <-input:
		kill <- true
		t.Error("Received Error:", logString)
	case <-time.After(time.Second * 5):
	}

	kill <- true

	if retVal != expected {
		t.Error("Expected:", expected, "Got:", retVal)
	}

	csvImporterCleanup("test")
}

func TestCSVToPersonBadNumFields(t *testing.T) {
	csvImporterSetup("test")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := ResultContext{
		"Incorrect Number of Fields | Expected: 5 | Got: 4",
		PersonContext{"test/test2.csv", 3, Person{0, PersonName{"", "", ""}, ""}},
	}
	input := make(chan string, 0)
	output := make(chan ResultContext, 0)
	kill := make(chan bool, 0)
	m := mock.NewMockProcessingContextInterface(ctrl)
	m.EXPECT().SetNumRecords("test/test2.csv", 2).Return(nil)

	go CSVImporter(input, output, m, input, kill)
	input <- "test/test2.csv"

	retVal := ResultContext{}
	for i := 0; i < 2; i++ {
		select {
		case retVal = <-output:
		case logString := <-input:
			kill <- true
			t.Error("Received Error:", logString)
		case <-time.After(time.Second * 5):
		}
	}

	kill <- true

	if retVal != expected {
		t.Error("Expected:", expected, "Got:", retVal)
	}

	csvImporterCleanup("test")
}

func TestCSVToPersonBadId(t *testing.T) {
	csvImporterSetup("test")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := ResultContext{
		"Unable to convert first field 'INTERNAL_ID' to an int | Got: a",
		PersonContext{"test/test3.csv", 3, Person{0, PersonName{"", "", ""}, ""}},
	}
	input := make(chan string, 0)
	output := make(chan ResultContext, 0)
	kill := make(chan bool, 0)
	m := mock.NewMockProcessingContextInterface(ctrl)
	m.EXPECT().SetNumRecords("test/test3.csv", 2).Return(nil)

	go CSVImporter(input, output, m, input, kill)
	input <- "test/test3.csv"

	retVal := ResultContext{}
	for i := 0; i < 2; i++ {
		select {
		case retVal = <-output:
		case logString := <-input:
			kill <- true
			t.Error("Received Error:", logString)
		case <-time.After(time.Second * 5):
		}
	}

	kill <- true

	if retVal != expected {
		t.Error("Expected:", expected, "Got:", retVal)
	}

	csvImporterCleanup("test")
}

func TestExportGetFileName(t *testing.T) {
	obj := Export{"test/test.ext", "dest", &ErrorFormat{}}
	expected := "dest/test.csv"
	actual := obj.GetFileName()

	if expected != actual {
		t.Error("Expected:", expected, "Got:", actual)
	}
}

func TestErrorFormatFirstAndTerminating(t *testing.T) {
	obj := ErrorFormat{1, "test"}
	expected := "LINE_NUM,ERROR_MSG\n1,test\n"

	actual := obj.FirstRecord() + obj.TerminatingString()

	if expected != actual {
		t.Error("Expected:", expected, "Got:", actual)
	}
}

func TestErrorFormatStandardRecord(t *testing.T) {
	obj := ErrorFormat{1, "test"}
	expected := "1,test\n"

	actual := obj.StandardRecord()

	if expected != actual {
		t.Error("Expected:", expected, "Got:", actual)
	}
}

func TestPersonFormatFirstAndTerminating(t *testing.T) {
	obj := PersonFormat{Person{1, PersonName{"a", "b", "c"}, "#"}}
	expected := `[{"id":1,"name":{"first":"a","middle":"b","last":"c"},"phone":"#"}]`
	actual := obj.FirstRecord() + obj.TerminatingString()
	if expected != actual {
		t.Error("Expected:", expected, "Got:", actual)
	}
}

func TestPersonFormatStandardRecord(t *testing.T) {
	obj := PersonFormat{Person{1, PersonName{"a", "b", "c"}, "#"}}
	expected := `,{"id":1,"name":{"first":"a","middle":"b","last":"c"},"phone":"#"}`
	actual := obj.StandardRecord()
	if expected != actual {
		t.Error("Expected:", expected, "Got:", actual)
	}
}

func TestPersonFormatNoMiddleName(t *testing.T) {
	obj := PersonFormat{Person{1, PersonName{"a", "", "c"}, "#"}}
	expected := `[{"id":1,"name":{"first":"a","last":"c"},"phone":"#"}]`
	actual := obj.FirstRecord() + obj.TerminatingString()
	if expected != actual {
		t.Error("Expected:", expected, "Got:", actual)
	}
}
