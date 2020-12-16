package main

import "errors"

// DistributeCSVRead splits the reading of a CSV file between a number of processes, as equitably
// as possible. Biggest challenge is maintaing line number count.

// DistributeWrite splits the writing of any file (first converted to []byte) between a number of
// processes, as equitably as possible.

// ReadCSVFile streams a CSV file from disk to the rest of the pipeline, record by record, as
// PersonContexts. When complete, it reports the total amount of records in the file.

// CSVToPerson converts a csv record from []string to Person.
func CSVToPerson(input []string) (Person, error) {
	return Person{}, errors.New("unimplemented")
}

// PersonToJSON converts a Person object to the JSON string representation.
func PersonToJSON(input Person) (string, error) {
	return "", errors.New("unimplemented")
}

// ErrorToCSV converts a line number and conversion error into a CSV string represenation.
func ErrorToCSV(lineNumber int, errString *string) (string, error) {
	return "", errors.New("unimplemented")
}
