package data

import "testing"

// Unit Test to verify validation when creatin a new Expense
func TestChecksValidation(t *testing.T) {
	e := &Expense{
		Vendor: "Super Market",
		Value:  78.20,
	}

	err := e.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
