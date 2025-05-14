package integration

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iMohamedSheta/xqb"
	"github.com/stretchr/testify/assert"
)

// setupTestDBForDelete creates a fresh test database connection and table with additional columns for delete tests
func setupTestDBForDelete(t *testing.T) *xqb.DBManager {
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
			email VARCHAR(255),
			age INT,
			status VARCHAR(50),
			price DECIMAL(10,2),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NULL
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

// insertTestDataForDelete inserts test data for delete operations
func insertTestDataForDelete(t *testing.T, dbManager *xqb.DBManager) {
	db, _ := dbManager.GetDB()
	_, err := db.Exec(`
		INSERT INTO test_users (name, email, age, status, price) VALUES
		('John Doe', 'john@example.com', 30, 'active', 100.00),
		('Jane Smith', 'jane@example.com', 25, 'active', 150.00),
		('Bob Wilson', 'bob@example.com', 35, 'inactive', 200.00)
	`)
	assert.NoError(t, err, "Failed to insert test data")
}

func TestDeleteBasic(t *testing.T) {
	dbManager := setupTestDBForDelete(t)
	defer cleanupTestDB(t, dbManager)

	tests := []struct {
		name     string
		where    map[string]interface{}
		wantRows int64
		wantErr  bool
		verifyFn func(*testing.T, *xqb.DBManager)
	}{
		{
			name: "delete single record",
			where: map[string]interface{}{
				"id": 1,
			},
			wantRows: 1,
			wantErr:  false,
			verifyFn: func(t *testing.T, dbManager *xqb.DBManager) {
				var count int
				db, _ := dbManager.GetDB()
				err := db.QueryRow("SELECT COUNT(*) FROM test_users WHERE id = 1").Scan(&count)
				assert.NoError(t, err)
				assert.Equal(t, 0, count)
			},
		},
		{
			name: "delete multiple records",
			where: map[string]interface{}{
				"status": "active",
			},
			wantRows: 2,
			wantErr:  false,
			verifyFn: func(t *testing.T, dbManager *xqb.DBManager) {
				var count int
				db, _ := dbManager.GetDB()
				err := db.QueryRow("SELECT COUNT(*) FROM test_users WHERE status = 'active'").Scan(&count)
				assert.NoError(t, err)
				assert.Equal(t, 0, count)
			},
		},
		{
			name: "delete with invalid condition",
			where: map[string]interface{}{
				"invalid_column": "value",
			},
			wantRows: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset test data before each test
			resetTestTable(t, dbManager)
			insertTestDataForDelete(t, dbManager)

			// Create a new query builder instance for each test
			qb := xqb.Table("test_users")

			// Build where conditions
			for col, val := range tt.where {
				qb.Where(col, "=", val)
			}

			// Perform delete
			got, err := qb.Delete()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantRows, got)

			// Run verification if provided
			if tt.verifyFn != nil {
				tt.verifyFn(t, dbManager)
			}
		})
	}
}

func TestDeleteWithTransaction(t *testing.T) {
	dbManager := setupTestDBForDelete(t)
	defer cleanupTestDB(t, dbManager)

	qb := xqb.Table("test_users")

	t.Run("failed transaction", func(t *testing.T) {
		testWithCleanTable(t, dbManager, func() {
			insertTestDataForDelete(t, dbManager)

			err := dbManager.Transaction(func(tx *sql.Tx) error {
				// First delete (success)
				_, err := qb.Where("id", "=", 1).DeleteTx(tx)
				assert.NoError(t, err)

				// Second delete with invalid condition
				_, err = qb.Where("invalid_column", "=", "value").DeleteTx(tx)
				assert.Error(t, err)

				return err // Return error to rollback transaction
			})
			assert.Error(t, err)

			// Verify no records were deleted (transaction was rolled back)
			var count int
			db, _ := dbManager.GetDB()
			err = db.QueryRow("SELECT COUNT(*) FROM test_users WHERE id = 1").Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 1, count) // Record should still exist
		})
	})
}
