package record

import (
	"errors"
	"fmt"
)

// A valid ID cannot be greater than 8 digits
const MAX_ID_INT = 99999999

// CsvHeader contains the structure of a valid csv header
// a valid csv header should contain the following names
// and exist in the given indexes
var CsvHeader = map[string]int{
	"INTERNAL_ID": 0,
	"FIRST_NAME":  1,
	"MIDDLE_NAME": 2,
	"LAST_NAME":   3,
	"PHONE_NUM":   4,
}

// Record contains the structure of a Record Object
type Record struct {
	Id       uint32      `json:"id"`
	Name     *Name       `json:"name"`
	PhoneNum PhoneNumber `json:"phone"`
}

// NewRecords returns a pointer to a record object
func NewRecord(id uint32, name *Name, phoneNum PhoneNumber) (*Record, error) {
	if id > MAX_ID_INT {
		return nil, errors.New(fmt.Sprint("ID exceeds the maximum length of ", MAX_ID_INT))
	}

	return &Record{
		Id:       id,
		Name:     name,
		PhoneNum: phoneNum,
	}, nil
}

// ErrorLog contains the structure of an ErrorLog object
type ErrorLog struct {
	Errors [][]string
}

// NewErrorLog returns a pointer to a new ErrorLog object
func NewErrorLog() *ErrorLog {
	return &ErrorLog{
		Errors: [][]string{{"LINE_NUM, ERROR_MSG"}},
	}
}

// Append takes a line number of an error and an error message
// and appends it to the current error log slice
func (el *ErrorLog) Append(lineNum string, msg string) {
	el.Errors = append(el.Errors, []string{lineNum, msg})
}

// HasData returns true if the ErrorLog has errors
func (el *ErrorLog) HasData() bool {
	return len(el.Errors) > 1
}
