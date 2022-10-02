package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (ex *Expenses) GetExpenses(rw http.ResponseWriter, r *http.Request) {
	ex.l.Println("Handle GET Expense")
	err := ex.db.GetExpenses(rw)
	if err != nil {
		http.Error(rw, "Could not retrieve database entries", http.StatusInternalServerError)
	}
}

func (ex *Expenses) GetExpense(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
	}
	ex.l.Println("Handle GET Specific Expense ", id)
	err = ex.db.GetExpense(rw, id)
	if err != nil {
		http.Error(rw, "Could not retrieve database entry", http.StatusInternalServerError)
	}
}

func (ex *Expenses) GetExpensesByVendor(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vendor, exists := vars["vendor"]

	if !exists {
		http.Error(rw, "Unable to read vendor", http.StatusBadRequest)
	}
	ex.l.Println("Handle GET Expenses By", vendor)
	err := ex.db.GetExpensesByVendor(rw, vendor)
	if err != nil {
		http.Error(rw, "Could not retrieve database entries", http.StatusInternalServerError)
	}
}
