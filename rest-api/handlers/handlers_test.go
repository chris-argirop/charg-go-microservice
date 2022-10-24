package handlers

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/chris-argirop/charg-go-microsrvice/rest-api/db"
	"github.com/gorilla/mux"
)

var jsonData = []byte(`{"vendor": "TestPAF","value": 12.3}`)

func Router() *mux.Router {
	router := mux.NewRouter()
	return router
}

func InitDB() *Expenses {
	l := log.New(os.Stdout, "charg-api", log.LstdFlags)

	dbaseDriver := os.Getenv("DB_DRIVER")
	dsourceName := os.Getenv("DS_NAME")

	db, err := db.NewDatabase(dbaseDriver, dsourceName)
	if err != nil {
		l.Fatal(err)
	}
	return NewExpense(l, db)
}

var testEh *Expenses

func TestMain(m *testing.M) {
	testEh = InitDB()
	defer testEh.db.CloseConnection()
	m.Run()
}

func TestAddExpenses(t *testing.T) {

	req, err := http.NewRequest("POST", "/api/add", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Errorf("Error creating a new HTTP Post request")
	}

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
	request, _ := http.NewRequest("GET", "/api/get", nil)
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(testEh.GetExpenses)

	handler.ServeHTTP(rec, request)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, expected %v", status, http.StatusOK)
	}
}

func TestUpdateExpenses(t *testing.T) {

	request, err := http.NewRequest("PUT", "/api/update/0", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}
	rec := httptest.NewRecorder()

	sm := Router()
	sm.HandleFunc("/api/update/{id:[0-9]+}", testEh.UpdateExpenses).Methods(http.MethodPut)
	sm.Use(testEh.MiddlewarevalidateExpense)
	sm.ServeHTTP(rec, request)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, expected %v", status, http.StatusOK)
	}

}
