package data

import (
	"encoding/json"
	"io"
)

type Month struct {
	Name        string
	TotalExp    float32
	AllExpenses []Expense
}

type Months []*Month

var Calendar = []*Month{
	&Month{
		Name:     "January",
		TotalExp: 0,
	},
	&Month{
		Name:     "February",
		TotalExp: 0,
	},
	&Month{
		Name:     "March",
		TotalExp: 0,
	},
	&Month{
		Name:     "April",
		TotalExp: 0,
	},
	&Month{
		Name:     "May",
		TotalExp: 0,
	},
	&Month{
		Name:     "June",
		TotalExp: 0,
	},
	&Month{
		Name:     "July",
		TotalExp: 0,
	},
	&Month{
		Name:     "August",
		TotalExp: 0,
	},
	&Month{
		Name:     "September",
		TotalExp: 0,
	},
	&Month{
		Name:     "October",
		TotalExp: 0,
	},
	&Month{
		Name:     "November",
		TotalExp: 0,
	},
	&Month{
		Name:     "December",
		TotalExp: 0,
	},
}

func CalculateMonthlyExpenses() {
	for _, exp := range ExpList {
		for _, mnth := range Calendar {
			if exp.Month == mnth.Name {
				mnth.TotalExp += exp.Value
				mnth.AllExpenses = append(mnth.AllExpenses, *exp)
			}
		}
	}
}

func GetCalendar() Months {
	CalculateMonthlyExpenses()
	return Calendar
}

func (cal *Months) FromJSON(w io.Reader) error {
	e := json.NewDecoder(w)
	return e.Decode(cal)
}

func (cal *Months) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(cal)
}
