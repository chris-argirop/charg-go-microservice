package handlers

import (
	"charg-go-microservice/rest-api/data"
	"log"
	"net/http"
)

type Expenses struct {
	l *log.Logger
}

func NewExpense(l *log.Logger) *Expenses {
	return &Expenses{l}
}

func (ex *Expenses) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ex.getExpenses(rw, r)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (ex *Expenses) getExpenses(rw http.ResponseWriter, r *http.Request) {
	le := data.GetExpenses()
	err := le.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to unmarshal JSON", http.StatusInternalServerError)
	}

}
