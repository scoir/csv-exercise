package main

import "testing"

func TestInternalIDPositive(t *testing.T) {
	expected := "| INTERNAL_ID must be positive and a maximum of 8 digits "
	input := Person{
		ID:    -1,
		Name:  PersonName{"a", "b", "c"},
		Phone: "111-111-1111",
	}

	result := VerifyRecord(&input)
	if result != expected {
		t.Error("Expected:", expected, "Got:", result)
	}
}

func TestInternalIDEightOrLessDigits(t *testing.T) {
	expected := "| INTERNAL_ID must be positive and a maximum of 8 digits "
	input := Person{
		ID:    100000000,
		Name:  PersonName{"a", "b", "c"},
		Phone: "111-111-1111",
	}

	result := VerifyRecord(&input)
	if result != expected {
		t.Error("Expected:", expected, "Got:", result)
	}
}

func TestFirstNameFifteenOrLessCharacters(t *testing.T) {
	expected := "| FIRST_NAME must be non-empty and 15 characters or less "
	input := Person{
		ID:    1,
		Name:  PersonName{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "b", "c"},
		Phone: "111-111-1111",
	}

	result := VerifyRecord(&input)
	if result != expected {
		t.Error("Expected:", expected, "Got:", result)
	}
}

func TestFirstNamePresent(t *testing.T) {
	expected := "| FIRST_NAME must be non-empty and 15 characters or less "
	input := Person{
		ID:    1,
		Name:  PersonName{"", "b", "c"},
		Phone: "111-111-1111",
	}

	result := VerifyRecord(&input)
	if result != expected {
		t.Error("Expected:", expected, "Got:", result)
	}
}

func TestMiddleNameFifteenOrLessCharacters(t *testing.T) {
	expected := "| MIDDLE_NAME must be 15 characters or less "
	input := Person{
		ID:    1,
		Name:  PersonName{"a", "bbbbbbbbbbbbbbbbbbbbb", "c"},
		Phone: "111-111-1111",
	}

	result := VerifyRecord(&input)
	if result != expected {
		t.Error("Expected:", expected, "Got:", result)
	}
}

func TestMiddleNameCanBeMissing(t *testing.T) {
	expected := ""
	input := Person{
		ID:    1,
		Name:  PersonName{"a", "", "c"},
		Phone: "111-111-1111",
	}

	result := VerifyRecord(&input)
	if result != expected {
		t.Error("Expected:", expected, "Got:", result)
	}
}

func TestLastNameFifteenOrLessCharacters(t *testing.T) {
	expected := "| LAST_NAME must be non-empty and 15 characters or less "
	input := Person{
		ID:    1,
		Name:  PersonName{"a", "b", "ccccccccccccccccccccccccc"},
		Phone: "111-111-1111",
	}

	result := VerifyRecord(&input)
	if result != expected {
		t.Error("Expected:", expected, "Got:", result)
	}
}

func TestLastNamePresent(t *testing.T) {
	expected := "| LAST_NAME must be non-empty and 15 characters or less "
	input := Person{
		ID:    1,
		Name:  PersonName{"a", "b", ""},
		Phone: "111-111-1111",
	}

	result := VerifyRecord(&input)
	if result != expected {
		t.Error("Expected:", expected, "Got:", result)
	}
}

func TestPhoneNumMatchesPattern(t *testing.T) {
	expected := "| phone number must be of the format ###-###-#### "
	input := Person{
		ID:    1,
		Name:  PersonName{"a", "b", "c"},
		Phone: "1a1-111-1111",
	}

	result := VerifyRecord(&input)
	if result != expected {
		t.Error("Expected:", expected, "Got:", result)
	}
}

func TestPhoneNumPresent(t *testing.T) {
	expected := "| phone number must be of the format ###-###-#### "
	input := Person{
		ID:    1,
		Name:  PersonName{"a", "b", "c"},
		Phone: "",
	}

	result := VerifyRecord(&input)
	if result != expected {
		t.Error("Expected:", expected, "Got:", result)
	}
}
