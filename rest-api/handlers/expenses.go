package handlers

import (
	"charg-go-microservice/rest-api/data"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Expenses struct {
	l *log.Logger
}

func NewExpense(l *log.Logger) *Expenses {
	return &Expenses{l}
}

func (ex *Expenses) GetExpenses(rw http.ResponseWriter, r *http.Request) {
	ex.l.Println("Handle GET Expense")
	le := data.GetExpenses()
	err := le.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to unmarshal JSON", http.StatusInternalServerError)
	}

}

func (ex *Expenses) AddExpenses(rw http.ResponseWriter, r *http.Request) {
	ex.l.Println("Handle POST Expense")

	exp := r.Context().Value(KeyExpense{}).(data.Expense)
	data.AddExpense(&exp)
}

type KeyExpense struct{}

func (e *Expenses) UpdateExpenses(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(rw, "Unable to conver id", http.StatusBadRequest)
	}
	e.l.Println("Handle PUT Expense", id)

	exp := r.Context().Value(KeyExpense{}).(data.Expense)

	err = data.UpdateExpense(id, &exp)
	if err != nil {
		http.Error(rw, "Expense not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Expense not found", http.StatusInternalServerError)
		return
	}
}

func (e *Expenses) MiddlewarevalidateExpense(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		exp := data.Expense{}

		err := exp.FromJSON(r.Body)
		if err != nil {
			e.l.Println("[ERROR] deserializing expense", err)
			http.Error(rw, "Error reading expense", http.StatusBadRequest)
			return
		}

		// validte the product
		err = exp.Validate()
		if err != nil {
			e.l.Println("[ERROR] validating expense", err)
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
