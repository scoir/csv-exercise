package main

import (
	"regexp"
)

// VerifyRecord returns "" if the Person object is valid, a string specifying the error(s)
// otherwise.
func VerifyRecord(person *Person) string {
	errMsg := ""
	// Check ID
	if !(person.ID > 0 && person.ID < 100000000) {
		errMsg += "| INTERNAL_ID must be positive and a maximum of 8 digits "
	}

	if result := verifyName(person.Name.First, "FIRST_NAME"); result != "" {
		errMsg += "| " + result + " "
	}

	if len(person.Name.Middle) > 15 {
		errMsg += "| MIDDLE_NAME must be 15 characters or less "
	}

	if result := verifyName(person.Name.Last, "LAST_NAME"); result != "" {
		errMsg += "| " + result + " "
	}

	result, mErr := regexp.MatchString(`\d\d\d-\d\d\d-\d\d\d\d`, person.Phone)
	if mErr != nil {
		errMsg += "| Unexpected regex match error: " + mErr.Error() + " "
	} else if !result {
		errMsg += "| phone number must be of the format ###-###-#### "
	}

	return errMsg
}

func verifyName(value, field string) string {
	if len(value) > 15 || value == "" {
		return field + " must be non-empty and 15 characters or less"
	}
	return ""
}

// PersonContext provides context to enable non-sequential verification of records.
type PersonContext struct {
	FileName   string
	LineNumber int // if the line number is "-1" it represents a "FILE" level error
	Record     Person
}

// ResultContext provides context to enable non-sequential reporting of record verification.
type ResultContext struct {
	result  string
	context PersonContext
}

// Verifier is a goroutine daemon to provide record verification without constant goroutine churn.
func Verifier(personChan chan PersonContext, resultChan chan ResultContext, killChan chan bool) {
	kill := false
	for !kill {
		select {
		case kill = <-killChan:
			continue
		case person := <-personChan:
			resultChan <- ResultContext{
				VerifyRecord(&person.Record),
				person,
			}
		}
	}
}
