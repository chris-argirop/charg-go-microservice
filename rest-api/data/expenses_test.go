package data

import "testing"

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
