package testhelpers

import (
	"database/sql"
	"os"
	"testing"

	"github.com/epalmerini/abitudini/internal/db"
)

// NewTestDB creates a temporary in-memory SQLite database for testing
func NewTestDB(t *testing.T) *sql.DB {
	t.Helper()

	// Create a temporary database file
	tmpFile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatalf("failed to create temp db file: %v", err)
	}
	tmpFile.Close()

	// Open and initialize the database
	testDB, err := db.Init(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to initialize test database: %v", err)
	}

	t.Cleanup(func() {
		testDB.Close()
		os.Remove(tmpFile.Name())
	})

	return testDB
}
