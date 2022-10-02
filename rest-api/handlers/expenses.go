package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/chris-argirop/charg-go-microsrvice/rest-api/data"
	"github.com/chris-argirop/charg-go-microsrvice/rest-api/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Expenses struct {
	l  *log.Logger
	db *db.Database
}

func NewExpense(l *log.Logger, db *db.Database) *Expenses {
	return &Expenses{l, db}
}

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

func (ex *Expenses) AddExpenses(rw http.ResponseWriter, r *http.Request) {
	ex.l.Println("Handle POST Expense")

	exp := r.Context().Value(KeyExpense{}).(data.Expense)
	err := ex.db.AddExpense(exp.Vendor, exp.Value)
	if err != nil {
		http.Error(rw, "Could not add database entry", http.StatusInternalServerError)
	}
}

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

func (ex *Expenses) DeleteExpense(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(rw, "Unable to conver id", http.StatusBadRequest)
	}
	ex.l.Println("Handle PUT Delete Specific Expense ", id)
	err = ex.db.DeleteExpense(id)
	if err != nil {
		http.Error(rw, "Could not remove database entry", http.StatusInternalServerError)
	}

}

func (ex *Expenses) ClearTable(rw http.ResponseWriter, r *http.Request) {
	ex.l.Println("Handle PUT Clear Table ")
	err := ex.db.ClearTable()

	if err != nil {
		http.Error(rw, "Could not clear table", http.StatusInternalServerError)
	}
}

func (ex *Expenses) MiddlewarevalidateExpense(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		exp := data.Expense{}

		err := exp.FromJSON(r.Body)
		if err != nil {
			ex.l.Println("[ERROR] deserializing expense", err)
			http.Error(rw, "Error reading expense", http.StatusBadRequest)
			return
		}

		// validte the product
		err = exp.Validate()
		if err != nil {
			ex.l.Println("[ERROR] validating expense", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating expense: %s", err),
				http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyExpense{}, exp)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
