package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Deletes an Expense based on the id provided in the GET Req path
func (ex *Expenses) DeleteExpense(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(rw, "Unable to conver id", http.StatusBadRequest)
	}
	ex.l.Println("Handle DELETE Specific Expense ", id)
	err = ex.db.DeleteExpense(id)
	if err != nil {
		http.Error(rw, "Could not remove database entry", http.StatusInternalServerError)
	}

}

// Completely removes any entries to the expenses table
func (ex *Expenses) ClearTable(rw http.ResponseWriter, r *http.Request) {
	ex.l.Println("Handle DELETE Clear Table ")
	err := ex.db.ClearTable()

	if err != nil {
		http.Error(rw, "Could not clear table", http.StatusInternalServerError)
	}
}
