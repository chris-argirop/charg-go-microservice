package handlers

import (
	"net/http"
	"strconv"

	"github.com/chris-argirop/charg-go-microsrvice/rest-api/data"
	"github.com/gorilla/mux"
)

type KeyExpense struct{}

func (ex *Expenses) UpdateExpenses(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(rw, "Unable to conver id", http.StatusBadRequest)
	}
	ex.l.Println("Handle PUT Expense", id)

	exp := r.Context().Value(KeyExpense{}).(data.Expense)

	//err = data.UpdateExpense(id, &exp)
	err = ex.db.UpdateExpense(id, exp.Vendor, exp.Description, exp.Value)
	if err != nil {
		http.Error(rw, "Expense not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Expense not found", http.StatusInternalServerError)
		return
	}
}
