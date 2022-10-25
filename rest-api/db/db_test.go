package db

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/go-playground/assert.v1"
)

var dbaseDriver string = os.Getenv("DB_DRIVER")
var dsourceName string = os.Getenv("DS_NAME")

func TestNewDatabase(t *testing.T) {

	_, err := NewDatabase(dbaseDriver, dsourceName)
	if err != nil {
		t.Fatalf("Creating New Database Connection: got %v, expected nil", err)
	}

}

func TestNewRow(t *testing.T) {
	testR := NewRow(1, "testVendor", "testDescr", 0.0, time.Now().String())
	require.NotEmpty(t, testR)
}

func TestToJSON(t *testing.T) {
	var jsonData = []byte(`{"vendor": "TestPAF","value": 12.3}`)
	testR := NewRow(1, "testVendor", "testDescr", 0.0, time.Now().String())

	err := testR.ToJSON(bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("JSON encoding: got %v, expected nil", err)
	}
}

func TestGetNextID(t *testing.T) {
	testDB, err := NewDatabase(dbaseDriver, dsourceName)
	require.NoError(t, err)

	// Clear DB
	err = testDB.ClearTable()
	require.NoError(t, err)
	// Retrieve the next ID
	nextID, err := testDB.getnextID()
	assert.Equal(t, nextID, 0)
	require.NoError(t, err)

	// Add an entry and re-calculate the next id
	err = testDB.AddExpense("testVendor", 0.0)
	require.NoError(t, err)
	nextID, err = testDB.getnextID()
	assert.Equal(t, nextID, 1)
	require.NoError(t, err)

	testDB.CloseConnection()
	_, err = testDB.getnextID()
	require.Error(t, err)

}

func TestDBExpenses(t *testing.T) {
	data := Row{}
	rw := httptest.NewRecorder()
	rw.Header().Set("Content-Type", "application/json")
	// Connect to the DB
	testDB, err := NewDatabase(dbaseDriver, dsourceName)
	defer testDB.CloseConnection()
	require.NoError(t, err)
	// Clear Table in DB
	clearErr := testDB.ClearTable()
	require.NoError(t, clearErr)
	// Add Expense to DB
	addErr := testDB.AddExpense("testVendor", 0.0)
	require.NoError(t, addErr)
	// Get the DB Entries
	err = testDB.GetExpenses(rw)
	require.NoError(t, err)
	// Verify that the entry added matches the input data
	json.Unmarshal(rw.Body.Bytes(), &data)
	assert.Equal(t, data.Vendor, "testVendor")
	assert.Equal(t, float64(data.Value), 0.0)

	// Update the first Entry
	err = testDB.UpdateExpense(0, "testVendor2", "N/A", 10.0)
	require.NoError(t, err)

	// Verify that the entry with ID 0 was updated
	rw = httptest.NewRecorder()
	err = testDB.GetExpense(rw, 0)
	require.NoError(t, err)
	json.Unmarshal(rw.Body.Bytes(), &data)
	assert.Equal(t, data.Vendor, "testVendor2")
	assert.Equal(t, float64(data.Value), 10.0)

	// Verify that retrieving an expense by a specific vendor is the correct one
	rw = httptest.NewRecorder()
	err = testDB.GetExpensesByVendor(rw, "testVendor2")
	require.NoError(t, err)
	json.Unmarshal(rw.Body.Bytes(), &data)
	assert.Equal(t, data.ID, 0)
	assert.Equal(t, float64(data.Value), 10.0)

	// Delete specific entry by ID
	err = testDB.DeleteExpense(0)
	require.NoError(t, err)

	// Verify Deletion
	err = testDB.GetExpense(rw, 0)
	require.NoError(t, err)

	json.Unmarshal(rw.Body.Bytes(), &data)
	require.NotEmpty(t, data)

}

func TestNoDB(t *testing.T) {
	testDB, err := NewDatabase(dbaseDriver, dsourceName)
	require.NoError(t, err)
	testDB.CloseConnection()

	err = testDB.AddExpense("testVendor", 0.0)
	require.Error(t, err)

	err = testDB.UpdateExpense(0, "testVendor2", "N/A", 10.0)
	require.Error(t, err)

	err = testDB.DeleteExpense(0)
	require.Error(t, err)
}
