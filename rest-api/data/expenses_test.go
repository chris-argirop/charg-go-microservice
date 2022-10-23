package data

import (
	"bytes"
	"io"
	"testing"
)

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

func TestFromJSON(t *testing.T) {
	e := &Expense{}
	var w io.Reader = bytes.NewBufferString(`{"vendor": "PAF", "value": 12.3}`)
	err := e.FromJSON(w)
	if err != nil {
		t.Fatalf("Decoding JSON FAILED")
	}

}
