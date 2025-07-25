package integration

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iMohamedSheta/xqb"
	"github.com/stretchr/testify/assert"
)

// setupTestDBForInsert creates a fresh test database connection and table
func setupTestDBForInsert(t *testing.T) *xqb.DBManager {
	db, err := dbManager.GetDB()
	if err != nil {
		t.Fatalf("Failed to get database connection: %v", err)
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

	return dbManager
}

// resetTestTable truncates the test table to ensure a clean state
func resetTestTableForInsert(t *testing.T, dbManager *xqb.DBManager) {
	db, _ := dbManager.GetDB()
	_, err := db.Exec("TRUNCATE TABLE test_users")
	assert.NoError(t, err, "Failed to reset test table")
}

// testWithCleanTable is a helper function that runs a test with a clean table
func testWithCleanTable(t *testing.T, dbManager *xqb.DBManager, testFn func()) {
	resetTestTableForInsert(t, dbManager)
	testFn()
}

func TestInsert(t *testing.T) {
	dbManager := setupTestDBForInsert(t)

	qb := xqb.Table("test_users")

	tests := []struct {
		name     string
		data     []map[string]any
		wantRows int64
		wantErr  bool
	}{
		{
			name: "single insert",
			data: []map[string]any{
				{"name": "John", "age": 30},
			},
			wantRows: 1,
			wantErr:  false,
		},
		{
			name: "multiple insert",
			data: []map[string]any{
				{"name": "John", "age": 30},
				{"name": "Jane", "age": 25},
			},
			wantRows: 2,
			wantErr:  false,
		},
		{
			name: "invalid data",
			data: []map[string]any{
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
				db, _ := dbManager.GetDB()
				err = db.QueryRow("SELECT COUNT(*) FROM test_users").Scan(&count)
				assert.NoError(t, err)
				assert.Equal(t, int(tt.wantRows), count)
			})
		})
	}
}

func TestInsertGetId(t *testing.T) {
	dbManager := setupTestDBForInsert(t)

	qb := xqb.Table("test_users")

	tests := []struct {
		name    string
		data    []map[string]any
		wantId  int64
		wantErr bool
	}{
		{
			name: "single insert",
			data: []map[string]any{
				{"name": "John", "age": 30},
			},
			wantId:  1, // First insert should have Id 1
			wantErr: false,
		},
		{
			name: "multiple insert",
			data: []map[string]any{
				{"name": "John", "age": 30},
				{"name": "Jane", "age": 25},
			},
			wantId:  1, // First insert in this test should have Id 1
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
	dbManager := setupTestDBForInsert(t)

	qb := xqb.Table("test_users")

	t.Run("failed transaction", func(t *testing.T) {
		testWithCleanTable(t, dbManager, func() {
			err := dbManager.Transaction(func(tx *sql.Tx) error {
				// First insert (success)
				_, err := qb.InsertTx([]map[string]any{
					{"name": "John", "age": 30},
				}, tx)
				assert.NoError(t, err)

				// Second insert (fail)
				_, err = qb.InsertTx([]map[string]any{
					{"invalid_column": "value"},
				}, tx)
				assert.Error(t, err)

				return err // Return error to rollback transaction
			})
			assert.Error(t, err)

			// Verify no records were inserted (transaction was rolled back)
			var count int
			db, _ := dbManager.GetDB()
			err = db.QueryRow("SELECT COUNT(*) FROM test_users").Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 0, count)
		})
	})

	t.Run("successful transaction", func(t *testing.T) {
		testWithCleanTable(t, dbManager, func() {
			err := dbManager.Transaction(func(tx *sql.Tx) error {
				// First insert
				affected, err := qb.InsertTx([]map[string]any{
					{"name": "John", "age": 30},
				}, tx)
				assert.NoError(t, err)
				assert.Equal(t, int64(1), affected)

				// Second insert
				affected, err = qb.InsertTx([]map[string]any{
					{"name": "Jane", "age": 25},
				}, tx)
				assert.NoError(t, err)
				assert.Equal(t, int64(1), affected)

				return nil
			})
			assert.NoError(t, err)

			// Verify both records were inserted
			var count int
			db, _ := dbManager.GetDB()
			err = db.QueryRow("SELECT COUNT(*) FROM test_users").Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 2, count)
		})
	})
}

func TestInsertGetIdWithTransaction(t *testing.T) {
	dbManager := setupTestDBForInsert(t)

	qb := xqb.Table("test_users")

	t.Run("failed transaction", func(t *testing.T) {
		testWithCleanTable(t, dbManager, func() {
			err := dbManager.Transaction(func(tx *sql.Tx) error {
				// First insert (success)
				_, err := qb.InsertGetIdTx([]map[string]any{
					{"name": "John", "age": 30},
				}, tx)
				assert.NoError(t, err)

				// Second insert (fail)
				_, err = qb.InsertGetIdTx([]map[string]any{
					{"invalid_column": "value"},
				}, tx)
				assert.Error(t, err)

				return err // Return error to rollback transaction
			})
			assert.Error(t, err)

			// Verify no records were inserted (transaction was rolled back)
			var count int
			db, _ := dbManager.GetDB()
			err = db.QueryRow("SELECT COUNT(*) FROM test_users").Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 0, count)
		})
	})

	t.Run("successful transaction", func(t *testing.T) {
		testWithCleanTable(t, dbManager, func() {
			err := dbManager.Transaction(func(tx *sql.Tx) error {
				// First insert
				id, err := qb.InsertGetIdTx([]map[string]any{
					{"name": "John", "age": 30},
				}, tx)
				assert.NoError(t, err)
				assert.Equal(t, int64(1), id)

				// Second insert
				id, err = qb.InsertGetIdTx([]map[string]any{
					{"name": "Jane", "age": 25},
				}, tx)
				assert.NoError(t, err)
				assert.Equal(t, int64(2), id)

				return nil
			})
			assert.NoError(t, err)

			// Verify both records were inserted
			var count int
			db, _ := dbManager.GetDB()
			err = db.QueryRow("SELECT COUNT(*) FROM test_users").Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 2, count)
		})
	})
}
