package main

// Person contains the complete record strcture for an identity.
type Person struct {
	ID    int        `json:"id"`
	Name  PersonName `json:"name"`
	Phone string     `json:"phone"`
}

// PersonName contains the fields which make up a Name for a Person.
type PersonName struct {
	First  string `json:"first"`
	Middle string `json:"middle,omitempty"`
	Last   string `json:"last"`
}
