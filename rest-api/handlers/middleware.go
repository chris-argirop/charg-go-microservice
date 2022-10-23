package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/chris-argirop/charg-go-microsrvice/rest-api/data"
)

// MiddleWare function to validate provided JSON in a request
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
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}
