package record

import (
	"errors"
	"fmt"
	"regexp"
)

type PhoneNumber string

// NewPhoneNumber returns a string representation of a PhoneNumber
// and validates that it follows the correct format of ###-###-####
func NewPhoneNumber(pn string) (PhoneNumber, error) {
	match, _ := regexp.MatchString("^[0-9]\\d{2}-\\d{3}-\\d{4}$", pn)
	if !match {
		fmt.Println("Invalid phone number:", pn)
		return "", errors.New("invalid phone number provided. Phone number expected to be in the following format: ###-###-####")
	}
	return PhoneNumber(pn), nil
}
