package main

import "errors"

// VerifyRecord returns "" if the Person object is valid, a string specifying the error(s)
// otherwise.
func VerifyRecord(person *Person) (string, error) {
	return "", errors.New("unimplemented")
}

// PersonContext provides context to enable non-sequential verification of records.
type PersonContext struct {
	FileName   string
	LineNumber int
	Record     *Person
}

// ResultContext provides context to enable non-sequential reporting of record verification.
type ResultContext struct {
	result  string
	context PersonContext
}

// Verifier is a goroutine daemon to provide record verification without constant goroutine churn.
func Verifier(personChan chan PersonContext, resultChan chan ResultContext) error {
	return errors.New("unimplemented")
}
