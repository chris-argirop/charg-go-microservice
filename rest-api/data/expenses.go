package data

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator"
)

type Expense struct {
	ID          int     `json:"id"`
	Vendor      string  `json:"vendor" validate:"required"`
	Description string  `json:"description"`
	Value       float32 `json:"value" validate:"required"`
	CreatedOn   string  `json:-`
	UpdatedOn   string  `json:-`
	DeletedOn   string  `json:-`
	Month       string  `json:-`
}

func (ex *Expense) FromJSON(w io.Reader) error {
	e := json.NewDecoder(w)
	return e.Decode(ex)
}

func (ex *Expense) Validate() error {
	validate := validator.New()
	return validate.Struct(ex)
}
