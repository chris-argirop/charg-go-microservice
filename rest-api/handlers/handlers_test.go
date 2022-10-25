package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

var jsonData = []byte(`{"vendor": "TestPAF","value": 12.3}`)

func TestAddExpenses(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/add", bytes.NewBuffer(jsonData))
	require.NoError(t, err)

	sm := Router()
	sm.HandleFunc("/api/add", testEh.AddExpenses).Methods(http.MethodPost)
	sm.Use(testEh.MiddlewarevalidateExpense)
	rec := httptest.NewRecorder()

	sm.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, expected %v", status, http.StatusOK)
	}

}

func TestGetExpenses(t *testing.T) {
	//Test All Get Requests
	sm := Router()
	sm.HandleFunc("/api/get", testEh.GetExpenses).Methods(http.MethodGet)
	sm.HandleFunc("/api/get/{id:[0-9]+}", testEh.GetExpense).Methods(http.MethodGet)
	sm.HandleFunc("/api/get/{vendor:(?s).*}", testEh.GetExpensesByVendor).Methods(http.MethodGet)

	// Gell All Expenses
	request, err := http.NewRequest("GET", "/api/get", nil)
	rec := httptest.NewRecorder()
	require.NoError(t, err)

	sm.ServeHTTP(rec, request)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, expected %v", status, http.StatusOK)
	}

	// Get Specific Expense by ID
	request, err = http.NewRequest("GET", "/api/get/0", nil)
	rec = httptest.NewRecorder()
	require.NoError(t, err)

	sm.ServeHTTP(rec, request)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, expected %v", status, http.StatusOK)
	}

	// Get Expensies that match the asked Vendor
	request, err = http.NewRequest("GET", "/api/get/testVendor1", nil)
	rec = httptest.NewRecorder()
	require.NoError(t, err)

	sm.ServeHTTP(rec, request)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, expected %v", status, http.StatusOK)
	}

}

func TestUpdateExpenses(t *testing.T) {

	request, err := http.NewRequest("PUT", "/api/update/0", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	rec := httptest.NewRecorder()

	sm := Router()
	sm.HandleFunc("/api/update/{id:[0-9]+}", testEh.UpdateExpenses).Methods(http.MethodPut)
	sm.Use(testEh.MiddlewarevalidateExpense)
	sm.ServeHTTP(rec, request)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, expected %v", status, http.StatusOK)
	}

}

func TestDeleteExpenses(t *testing.T) {

	sm := Router()
	sm.HandleFunc("/api/delete/{id:[0-9]+}", testEh.DeleteExpense).Methods(http.MethodDelete)
	sm.HandleFunc("/api/clearall", testEh.ClearTable).Methods(http.MethodDelete)

	request, err := http.NewRequest("DELETE", "/api/delete/0", nil)
	if err != nil {
		require.NoError(t, err)
	}
	rec := httptest.NewRecorder()

	sm.ServeHTTP(rec, request)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, expected %v", status, http.StatusOK)
	}

	request, err = http.NewRequest("DELETE", "/api/clearall", nil)
	if err != nil {
		require.NoError(t, err)
	}
	rec = httptest.NewRecorder()

	sm.ServeHTTP(rec, request)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, expected %v", status, http.StatusOK)
	}

}
