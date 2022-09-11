package db

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	db *sql.DB
}

type Row struct {
	ID          int     `json:"id" validate:"required"`
	Vendor      string  `json:"vendor" validate:"required"`
	Description string  `json:"description"`
	Value       float32 `json:"value" validate:"required"`
	Date        string  `json:"createdOn"`
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

func NewRow(id int, vendor string, descr string, value float32, date string) *Row {
	return &Row{id, vendor, descr, value, date}
}

func (mydb Database) CloseConnection() {
	err := mydb.db.Close()
	if err != nil {
		log.Fatalf("DB Close Failure: %v", err)
	}
}

func (r Row) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(r)
}

func (mydb Database) AddExpense(vendor string, sum float32) error {
	id := mydb.getnextID()
	queryString, args, err := sq.Insert("expenses").Columns("id", "vendor", "descr", "val", "createdOn").
		Values(id, vendor, "N/A", sum, time.Now()).ToSql()

	if err != nil {
		return err
	}

	_, err = mydb.db.Exec(queryString, args...)
	if err != nil {
		return err
	}

	return nil
}

func (mydb Database) GetExpenses(w io.Writer) error {
	res, err := mydb.db.Query("SELECT * FROM expenses")
	if err != nil {
		log.Fatalf("DB Get Failure: %v", err)
		return err
	}

	defer res.Close()
	e := json.NewEncoder(w)
	for res.Next() {
		var (
			id     int
			vendor string
			descr  string
			value  float32
			date   string
		)

		err = res.Scan(&id, &vendor, &descr, &value, &date)
		if err != nil {
			log.Fatalf("DB Row parse Failure: %v", err)
		}
		err = e.Encode(NewRow(id, vendor, descr, value, date))
		if err != nil {
			log.Fatalf("JSON Encoding failed")
		}
	}

	return nil
}

func (mydb *Database) getnextID() int {
	var count int
	err := mydb.db.QueryRow("SELECT COUNT(*) FROM expenses").Scan(&count)
	if err != nil {
		log.Fatalf("Couldn't retrieve number of entries: %v", err)
	}
	return count
}
