package handlers

import (
	"log"
	"os"
	"testing"

	"github.com/chris-argirop/charg-go-microsrvice/rest-api/db"
	"github.com/gorilla/mux"
)

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
	os.Exit(m.Run())

	defer testEh.db.CloseConnection()
}
