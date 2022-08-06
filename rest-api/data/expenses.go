package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/go-playground/validator/v10"
)

type Expense struct {
	ID          int     `json:"id"`
	Vendor      string  `json:"vendor" validate:"required"`
	Description string  `json:"description"`
	Value       float32 `json:"value" validate:"required"`
	CreatedOn   string  `json:-`
	UpdatedOn   string  `json:-`
	DeletedOn   string  `json:-`
}

type Expenses []*Expense

func (ex *Expense) FromJSON(w io.Reader) error {
	e := json.NewDecoder(w)
	return e.Decode(ex)
}

func (ex *Expenses) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(ex)
}

func (e *Expense) Validate() error {
	validate := validator.New()
	return validate.Struct(e)
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

func GetExpenses() Expenses {
	return ExpList
}

func AddExpense(e *Expense) {
	e.ID = getnextID()
	ExpList = append(ExpList, e)
}

func UpdateExpense(id int, e *Expense) error {
	_, pos, err := findExpense(id)
	if err != nil {
		return err
	}

	e.ID = id
	ExpList[pos] = e
	return nil
}

var ErrExpenseNotFound = fmt.Errorf("Expense not found!")

func findExpense(id int) (*Expense, int, error) {
	for i, e := range ExpList {
		if e.ID == id {
			return e, i, nil
		}

	}
	return nil, -1, ErrExpenseNotFound
}

func getnextID() int {
	le := ExpList[len(ExpList)-1]
	return le.ID + 1
}
