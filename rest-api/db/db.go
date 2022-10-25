package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

// Database Constructor
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

// Row Constructor
func NewRow(id int, vendor string, descr string, value float32, date string) *Row {
	return &Row{id, vendor, descr, value, date}
}

// Close connection to the DB
func (mydb Database) CloseConnection() {
	err := mydb.db.Close()
	if err != nil {
		log.Fatalf("DB Close Failure: %v", err)
	}
}

// JSON Encoder Function from an io.Writer
func (r Row) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(r)
}

// Write on a io.Writer all the Expense entries from the DB
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

// Write on a io.Writer the Expense entry from the DB, by id
func (mydb Database) GetExpense(w io.Writer, id int) error {
	res, err := mydb.db.Query(fmt.Sprintf("SELECT * FROM expenses WHERE ID = %d", id))
	if err != nil {
		log.Fatalf("DB Get Single Expense Failure: %v", err)
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
			log.Fatalf("DB Single Row parse Failure: %v", err)
		}
		err = e.Encode(NewRow(id, vendor, descr, value, date))
		if err != nil {
			log.Fatalf("JSON Encoding failed")
		}
	}

	return nil
}

// Write on a io.Writer all the Expense entries from the DB, that match a specific vendor
func (mydb Database) GetExpensesByVendor(w io.Writer, vendor string) error {
	res, err := mydb.db.Query(fmt.Sprintf("SELECT * FROM expenses WHERE vendor = \"%s\"", vendor))
	if err != nil {
		log.Fatalf("DB Get Expenses by Vendor Failure: %v", err)
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
			log.Fatalf("DB Single Row parse Failure: %v", err)
		}
		err = e.Encode(NewRow(id, vendor, descr, value, date))
		if err != nil {
			log.Fatalf("JSON Encoding failed")
		}
	}

	return nil
}

// Add a new expense entry in the DB, information needed to be supplied as paramaeters
// vendor: a string describing the vendor that the expense was made to
// sum: a float32 to specify the ammount of the expense to the above vendor
func (mydb Database) AddExpense(vendor string, sum float32) error {
	id, err := mydb.getnextID()
	if err != nil {
		return err
	}
	queryString, args, err := sq.Insert("expenses").
		Columns("id", "vendor", "descr", "val", "createdOn").
		Values(id, vendor, "N/A", sum, time.Now()).
		ToSql()

	if err != nil {
		return err
	}

	_, err = mydb.db.Exec(queryString, args...)
	if err != nil {
		return err
	}

	return nil
}

// Update a DB entry that matches a specific id, with the values provided in the parameters
func (mydb Database) UpdateExpense(id int, vendor string, description string, sum float32) error {

	queryString, args, err := sq.Update("expenses").
		Set("vendor", vendor).
		Set("descr", description).
		Set("val", sum).
		Set("createdOn", time.Now()).
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return err
	}

	_, err = mydb.db.Exec(queryString, args...)
	if err != nil {
		return err
	}

	return nil

}

// Delete a specific DB entry that matches the supplied id value
func (mydb Database) DeleteExpense(id int) error {
	queryString, args, err := sq.Delete("expenses").Where("id = ?", id).ToSql()

	if err != nil {
		return err
	}
	_, err = mydb.db.Exec(queryString, args...)
	if err != nil {
		return err
	}

	return nil
}

// Clear all entries from the DB table
func (mydb Database) ClearTable() error {
	queryString, args, err := sq.Delete("expenses").ToSql()

	if err != nil {
		return err
	}

	_, err = mydb.db.Exec(queryString, args...)
	if err != nil {
		return err
	}

	return nil
}

// Function to retrieve the ID which a new db entry should have
func (mydb *Database) getnextID() (int, error) {
	var count int
	err := mydb.db.QueryRow("SELECT COUNT(*) FROM expenses").Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}
