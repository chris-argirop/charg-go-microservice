package data

import (
	"encoding/json"
	"io"
	"time"
)

type Expense struct {
	ID          int     `json:"id"`
	Vendor      string  `json:"vendor"`
	Description string  `json:"description"`
	Value       float32 `json:"value"`
	CreatedOn   string  `json:-`
	UpdatedOn   string  `json:-`
	DeletedOn   string  `json:-`
}

type Expenses []*Expense

func (ex *Expenses) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(ex)
}

func GetExpenses() Expenses {
	return ExpList
}

var ExpList = []*Expense{
	&Expense{
		ID:          1,
		Vendor:      "supermarket",
		Description: "Groceries",
		Value:       60.0,
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Expense{
		ID:          2,
		Vendor:      "efood",
		Description: "Take-out",
		Value:       11.90,
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
