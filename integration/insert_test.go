package main

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iMohamedSheta/xqb"
	"github.com/stretchr/testify/assert"
)

// setupTestDB creates a fresh test database connection and table
func setupTestDB(t *testing.T) *xqb.DBManager {
	// Create a test database connection
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/test_xqb_db")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Drop and recreate test table to ensure clean state
	_, err = db.Exec("DROP TABLE IF EXISTS test_users")
	if err != nil {
		t.Fatalf("Failed to drop test table: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE test_users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255),
			age INT
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Set up DBManager
	dbManager := xqb.GetDBManager()
	dbManager.SetDB(db)

	return dbManager
}

// cleanupTestDB closes the database connection
func cleanupTestDB(t *testing.T, dbManager *xqb.DBManager) {
	if err := dbManager.CloseDB(); err != nil {
		t.Errorf("Failed to close database: %v", err)
	}
}

// resetTestTable truncates the test table to ensure a clean state
func resetTestTable(t *testing.T, dbManager *xqb.DBManager) {
	_, err := dbManager.db.Exec("TRUNCATE TABLE test_users")
	assert.NoError(t, err, "Failed to reset test table")
}

// testWithCleanTable is a helper function that runs a test with a clean table
func testWithCleanTable(t *testing.T, dbManager *xqb.DBManager, testFn func()) {
	resetTestTable(t, dbManager)
	testFn()
}

func TestInsert(t *testing.T) {
	dbManager := setupTestDB(t)
	defer cleanupTestDB(t, dbManager)

	qb := xqb.Table("test_users")

	tests := []struct {
		name     string
		data     []map[string]interface{}
		wantRows int64
		wantErr  bool
	}{
		{
			name: "single insert",
			data: []map[string]interface{}{
				{"name": "John", "age": 30},
			},
			wantRows: 1,
			wantErr:  false,
		},
		{
			name: "multiple insert",
			data: []map[string]interface{}{
				{"name": "John", "age": 30},
				{"name": "Jane", "age": 25},
			},
			wantRows: 2,
			wantErr:  false,
		},
		{
			name: "invalid data",
			data: []map[string]interface{}{
				{"invalid_column": "value"},
			},
			wantRows: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testWithCleanTable(t, dbManager, func() {
				got, err := qb.Insert(tt.data)
				if tt.wantErr {
					assert.Error(t, err)
					return
				}
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRows, got)

				// Verify data was actually inserted
				var count int
				err = dbManager.db.QueryRow("SELECT COUNT(*) FROM test_users").Scan(&count)
				assert.NoError(t, err)
				assert.Equal(t, int(tt.wantRows), count)
			})
		})
	}
}

func TestInsertGetId(t *testing.T) {
	dbManager := setupTestDB(t)
	defer cleanupTestDB(t, dbManager)

	qb := xqb.Table("test_users")

	tests := []struct {
		name    string
		data    []map[string]interface{}
		wantId  int64
		wantErr bool
	}{
		{
			name: "single insert",
			data: []map[string]interface{}{
				{"name": "John", "age": 30},
			},
			wantId:  1, // First insert should have ID 1
			wantErr: false,
		},
		{
			name: "multiple insert",
			data: []map[string]interface{}{
				{"name": "John", "age": 30},
				{"name": "Jane", "age": 25},
			},
			wantId:  1, // First insert in this test should have ID 1
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testWithCleanTable(t, dbManager, func() {
				got, err := qb.InsertGetId(tt.data)
				if tt.wantErr {
					assert.Error(t, err)
					return
				}
				assert.NoError(t, err)
				assert.Equal(t, tt.wantId, got)
			})
		})
	}
}

func TestInsertWithTransaction(t *testing.T) {
	dbManager := setupTestDB(t)
	defer cleanupTestDB(t, dbManager)

	qb := xqb.Table("test_users")

	t.Run("failed transaction", func(t *testing.T) {
		testWithCleanTable(t, dbManager, func() {
			err := dbManager.Transaction(func(tx *sql.Tx) error {
				// First insert (success)
				_, err := qb.InsertWithTx([]map[string]interface{}{
					{"name": "John", "age": 30},
				}, tx)
				assert.NoError(t, err)

				// Second insert (fail)
				_, err = qb.InsertWithTx([]map[string]interface{}{
					{"invalid_column": "value"},
				}, tx)
				assert.Error(t, err)

				return err // Return error to rollback transaction
			})
			assert.Error(t, err)

			// Verify no records were inserted (transaction was rolled back)
			var count int
			err = dbManager.db.QueryRow("SELECT COUNT(*) FROM test_users").Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 0, count)
		})
	})

	t.Run("successful transaction", func(t *testing.T) {
		testWithCleanTable(t, dbManager, func() {
			err := dbManager.Transaction(func(tx *sql.Tx) error {
				// First insert
				affected, err := qb.InsertWithTx([]map[string]interface{}{
					{"name": "John", "age": 30},
				}, tx)
				assert.NoError(t, err)
				assert.Equal(t, int64(1), affected)

				// Second insert
				affected, err = qb.InsertWithTx([]map[string]interface{}{
					{"name": "Jane", "age": 25},
				}, tx)
				assert.NoError(t, err)
				assert.Equal(t, int64(1), affected)

				return nil
			})
			assert.NoError(t, err)

			// Verify both records were inserted
			var count int
			err = dbManager.db.QueryRow("SELECT COUNT(*) FROM test_users").Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 2, count)
		})
	})
}

func TestInsertGetIdWithTransaction(t *testing.T) {
	dbManager := setupTestDB(t)
	defer cleanupTestDB(t, dbManager)

	qb := xqb.Table("test_users")

	t.Run("failed transaction", func(t *testing.T) {
		testWithCleanTable(t, dbManager, func() {
			err := dbManager.Transaction(func(tx *sql.Tx) error {
				// First insert (success)
				_, err := qb.InsertGetIdWithTx([]map[string]interface{}{
					{"name": "John", "age": 30},
				}, tx)
				assert.NoError(t, err)

				// Second insert (fail)
				_, err = qb.InsertGetIdWithTx([]map[string]interface{}{
					{"invalid_column": "value"},
				}, tx)
				assert.Error(t, err)

				return err // Return error to rollback transaction
			})
			assert.Error(t, err)

			// Verify no records were inserted (transaction was rolled back)
			var count int
			err = dbManager.db.QueryRow("SELECT COUNT(*) FROM test_users").Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 0, count)
		})
	})

	t.Run("successful transaction", func(t *testing.T) {
		testWithCleanTable(t, dbManager, func() {
			err := dbManager.Transaction(func(tx *sql.Tx) error {
				// First insert
				id, err := qb.InsertGetIdWithTx([]map[string]interface{}{
					{"name": "John", "age": 30},
				}, tx)
				assert.NoError(t, err)
				assert.Equal(t, int64(1), id)

				// Second insert
				id, err = qb.InsertGetIdWithTx([]map[string]interface{}{
					{"name": "Jane", "age": 25},
				}, tx)
				assert.NoError(t, err)
				assert.Equal(t, int64(2), id)

				return nil
			})
			assert.NoError(t, err)

			// Verify both records were inserted
			var count int
			err = dbManager.db.QueryRow("SELECT COUNT(*) FROM test_users").Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 2, count)
		})
	})
}
