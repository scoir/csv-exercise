package main

// Person contains the complete record strcture for an identity.
type Person struct {
	ID    int
	Name  PersonName
	Phone string
}

// PersonName contains the fields which make up a Name for a Person.
type PersonName struct {
	First, Middle, Last string
}
