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

func TestAddExpenses(t *testing.T) {
	data := Row{}
	rw := httptest.NewRecorder()
	rw.Header().Set("Content-Type", "application/json")
	testDB, err := NewDatabase(dbaseDriver, dsourceName)
	require.NoError(t, err)
	clearErr := testDB.ClearTable()
	require.NoError(t, clearErr)
	addErr := testDB.AddExpense("testVendor", 0.0)
	require.NoError(t, addErr)
	err = testDB.GetExpenses(rw)
	require.NoError(t, err)

	json.Unmarshal(rw.Body.Bytes(), &data)
	assert.Equal(t, data.Vendor, "testVendor")
	assert.Equal(t, float64(data.Value), 0.0)
}
