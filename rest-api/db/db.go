package db

import (
	"database/sql"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(driverName, dataSourceName string) (*Database, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatalf("DB Connection Failure: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("DB Ping Failure: %v", err)
	}

	return &Database{db: db}, err
}

func (mydb Database) CloseConnection() {
	err := mydb.db.Close()
	if err != nil {
		log.Fatalf("DB Close Failure: %v", err)
	}
}

func (mydb Database) AddToHistory(vendor string, sum float32) error {
	queryString, args, err := sq.Insert("expenses").Columns("vendor", "val", "createdOn").
		Values(vendor, sum, time.Now()).ToSql()

	if err != nil {
		return err
	}

	_, err = mydb.db.Exec(queryString, args...)
	if err != nil {
		return err
	}

	return nil
}
