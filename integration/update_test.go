package integration

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iMohamedSheta/xqb"
	"github.com/stretchr/testify/assert"
)

// setupTestDBForUpdate creates a fresh test database connection and table with additional columns for update tests
func setupTestDBForUpdate(t *testing.T) *xqb.DBManager {
	// Create a test database connection
	db, err := dbManager.GetDB()
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

	return dbManager
}

// insertTestDataForUpdate inserts test data for update operations
func insertTestDataForUpdate(t *testing.T, dbManager *xqb.DBManager) {
	db, _ := dbManager.GetDB()
	_, err := db.Exec(`
		INSERT INTO test_users (name, email, age, status, price) VALUES
		('John Doe', 'john@example.com', 30, 'active', 100.00),
		('Jane Smith', 'jane@example.com', 25, 'active', 150.00),
		('Bob Wilson', 'bob@example.com', 35, 'inactive', 200.00)
	`)
	assert.NoError(t, err, "Failed to insert test data")
}

func TestUpdateBasic(t *testing.T) {
	dbManager := setupTestDBForUpdate(t)

	insertTestDataForUpdate(t, dbManager)

	tests := []struct {
		name     string
		where    map[string]any
		updates  map[string]any
		wantRows int64
		wantErr  bool
		verifyFn func(*testing.T, *xqb.DBManager)
	}{
		{
			name: "update single record",
			where: map[string]any{
				"id": 1,
			},
			updates: map[string]any{
				"name":  "John Updated",
				"email": "john.updated@example.com",
			},
			wantRows: 1,
			wantErr:  false,
			verifyFn: func(t *testing.T, dbManager *xqb.DBManager) {
				var name, email string
				db, _ := dbManager.GetDB()
				err := db.QueryRow("SELECT name, email FROM test_users WHERE id = 1").Scan(&name, &email)
				assert.NoError(t, err)
				assert.Equal(t, "John Updated", name)
				assert.Equal(t, "john.updated@example.com", email)
			},
		},
		{
			name: "update multiple records",
			where: map[string]any{
				"status": "active",
			},
			updates: map[string]any{
				"status": "pending",
			},
			wantRows: 2,
			wantErr:  false,
			verifyFn: func(t *testing.T, dbManager *xqb.DBManager) {
				var count int
				db, _ := dbManager.GetDB()
				err := db.QueryRow("SELECT COUNT(*) FROM test_users WHERE status = 'pending'").Scan(&count)
				assert.NoError(t, err)
				assert.Equal(t, 2, count)
			},
		},
		{
			name: "update with invalid column",
			where: map[string]any{
				"id": 1,
			},
			updates: map[string]any{
				"invalid_column": "value",
			},
			wantRows: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := xqb.Table("test_users")
			// Reset test data before each test
			resetTestTableForInsert(t, dbManager)
			insertTestDataForUpdate(t, dbManager)

			// Build where conditions
			for col, val := range tt.where {
				qb.Where(col, "=", val)
			}

			// Perform update
			got, err := qb.Update(tt.updates)
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

func TestUpdateWithTransaction(t *testing.T) {
	dbManager := setupTestDBForUpdate(t)

	qb := xqb.Table("test_users")

	t.Run("failed transaction", func(t *testing.T) {
		testWithCleanTable(t, dbManager, func() {
			insertTestDataForUpdate(t, dbManager)

			err := dbManager.Transaction(func(tx *sql.Tx) error {
				// First update (success)
				_, err := qb.Where("id", "=", 1).UpdateTx(map[string]any{
					"name": "John Updated",
				}, tx)
				assert.NoError(t, err)

				// Second update (fail)
				_, err = qb.Where("id", "=", 2).UpdateTx(map[string]any{
					"invalid_column": "value",
				}, tx)
				assert.Error(t, err)

				return err // Return error to rollback transaction
			})
			assert.Error(t, err)

			// Verify no records were updated (transaction was rolled back)
			var name string
			db, _ := dbManager.GetDB()
			err = db.QueryRow("SELECT name FROM test_users WHERE id = 1").Scan(&name)
			assert.NoError(t, err)
			assert.Equal(t, "John Doe", name) // Original value should remain
		})
	})
}
