package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// DistributeCSVRead splits the reading of a CSV file between a number of processes, as equitably
// as possible. Need to solve maintaing line number count.

// DistributeWrite splits the writing of any file (first converted to []byte) between a number of
// processes, as equitably as possible.

// CSVImporter streams CSV files from disk to the rest of the pipeline, record by record, as
// ResultContexts. When each file completes, the total amount of records in the file are reported.
func CSVImporter(
	fileNameStream chan string,
	output chan ResultContext,
	processingContext ProcessingContextInterface,
	logger chan string,
	killChan chan bool,
) {
	kill := false
	for !kill {
		select {
		case kill = <-killChan:
			continue
		case fileName := <-fileNameStream:
			if err := readCSVFile(fileName, output, processingContext); err != nil {
				LogError(err, logger)
			}
		}
	}
}

// readCSVFile reads a single file to memory.
func readCSVFile(
	fileName string,
	output chan ResultContext,
	processingContext ProcessingContextInterface,
) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	lineNum := 1
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	firstLine := strings.TrimSpace(scanner.Text())
	expectedHeader := "INTERNAL_ID,FIRST_NAME,MIDDLE_NAME,LAST_NAME,PHONE_NUM"
	if firstLine != expectedHeader {
		output <- ResultContext{
			fmt.Sprintf("Incorrect CSV Header | Expected: %s | Got: %s", expectedHeader, firstLine),
			PersonContext{
				fileName,
				lineNum,
				Person{},
			},
		}
		return nil
	}
	for scanner.Scan() {
		lineNum++
		rawLine := strings.TrimSpace(scanner.Text())

		result := ""
		person, cErr := csvToPerson(rawLine)
		if cErr != nil {
			result = cErr.Error()
		}
		output <- ResultContext{
			result,
			PersonContext{
				fileName,
				lineNum,
				person,
			},
		}
	}

	if err := processingContext.SetNumRecords(fileName, lineNum-1); err != nil {
		return err
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// csvToPerson converts a csv record from string to Person.
func csvToPerson(input string) (Person, error) {
	fields := strings.Split(input, ",")
	if len(fields) != 5 {
		return Person{}, fmt.Errorf(
			"Incorrect Number of Fields | Expected: 5 | Got: %d",
			len(fields),
		)
	}
	id, err := strconv.Atoi(fields[0])
	if err != nil {
		return Person{}, fmt.Errorf(
			"Unable to convert first field 'INTERNAL_ID' to an int | Got: %s",
			fields[0],
		)
	}
	return Person{id, PersonName{fields[1], fields[2], fields[3]}, fields[4]}, nil
}

// FormatInterface provides a mock target for Format objects
type FormatInterface interface {
	GetFileExtension() string
	FirstRecord() string
	StandardRecord() string
	TerminatingString() string
}

// ErrorFormat implements FormatInterface for Error-specification exports
type ErrorFormat struct {
	lineNum int
	detail  string
}

// GetFileExtension returns the correct file extension for the file format
func (t *ErrorFormat) GetFileExtension() string {
	return ".csv"
}

// FirstRecord returns the correct format for the first record in the file
func (t *ErrorFormat) FirstRecord() string {
	header := "LINE_NUM,ERROR_MSG"
	return fmt.Sprintf("%s\n%s", header, t.StandardRecord())
}

// StandardRecord returns the correct format for normal records in the file
func (t *ErrorFormat) StandardRecord() string {
	return fmt.Sprintf("%d,%s\n", t.lineNum, t.detail)
}

// TerminatingString returns the final characters in the file
// NOTE: This cannot include ANY record-specific information
func (t *ErrorFormat) TerminatingString() string {
	return ""
}

// PersonFormat implements FormatInterface for Person-specification exports
type PersonFormat struct {
	person Person
}

// GetFileExtension returns the correct file extension for the file format
func (t *PersonFormat) GetFileExtension() string {
	return ".json"
}

// FirstRecord returns the correct format for the first record in the file
func (t *PersonFormat) FirstRecord() string {
	rawJSON, _ := json.Marshal(t.person)
	return fmt.Sprintf("[%s", rawJSON)
}

// StandardRecord returns the correct format for normal records in the file
func (t *PersonFormat) StandardRecord() string {
	rawJSON, _ := json.Marshal(t.person)
	return fmt.Sprintf(",%s", rawJSON)
}

// TerminatingString returns the final characters in the file
// NOTE: This cannot include ANY record-specific information
func (t *PersonFormat) TerminatingString() string {
	return "]"
}
