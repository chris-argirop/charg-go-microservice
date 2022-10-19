package handlers

import (
	"net/http"

	"github.com/chris-argirop/charg-go-microsrvice/rest-api/data"
)

// Add a new expense to the DB, take input from the data input in JSON format provided in the POST REQ
func (ex *Expenses) AddExpenses(rw http.ResponseWriter, r *http.Request) {
	ex.l.Println("Handle POST Expense")

	exp := r.Context().Value(KeyExpense{}).(data.Expense)
	err := ex.db.AddExpense(exp.Vendor, exp.Value)
	if err != nil {
		http.Error(rw, "Could not add database entry", http.StatusInternalServerError)
	}
}
