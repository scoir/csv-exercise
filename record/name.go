package record

import (
	"errors"
	"fmt"
)

// Name stores the structure for a Name object
type Name struct {
	FirstName  string `json:"first"`
	MiddleName string `json:"middle,omitempty"`
	LastName   string `json:"last"`
}

// NewName returns a pointer to a new Name object
func NewName(first string, middle string, last string) (*Name, error) {
	err := validateName(first)

	if err != nil {
		return nil, errors.New(fmt.Sprint("Error in first name:", err))
	}

	err = validateName(last)

	if err != nil {
		return nil, errors.New(fmt.Sprint("error in last name:", err))
	}

	return &Name{
		FirstName:  first,
		MiddleName: middle,
		LastName:   last,
	}, nil
}

// validateName ensures a given name follows the correct naming rules
// and returns an error if the name violates the naming rules
func validateName(name string) error {
	if len(name) > 15 {
		return errors.New("Name exceeds maximum length of 15 characters")
	}

	if len(name) == 0 {
		return errors.New("Name is empty")
	}

	return nil
}
