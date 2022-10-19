package handlers

import (
	"log"

	"github.com/chris-argirop/charg-go-microsrvice/rest-api/db"
	_ "github.com/go-sql-driver/mysql"
)

type Expenses struct {
	l  *log.Logger
	db *db.Database
}

// Expenses Constructor
func NewExpense(l *log.Logger, db *db.Database) *Expenses {
	return &Expenses{l, db}
}
