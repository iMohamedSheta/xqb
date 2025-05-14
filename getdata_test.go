package xqb

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

// setupTestDB creates a fresh test database connection and table
// setupTestDBForUpdate creates a fresh test database connection and table with additional columns for update tests
func setupTestDB(t *testing.T) *DBManager {
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
	dbManager := GetDBManager()
	dbManager.SetDB(db)

	return dbManager
}

// insertTestDataForUpdate inserts test data for update operations
func insertTestDataForTest(t *testing.T, dbManager *DBManager) {
	db, _ := dbManager.GetDB()
	_, err := db.Exec(`
		INSERT INTO test_users (name, email, age, status, price) VALUES
		('John Doe', 'john@example.com', 30, 'active', 100.00),
		('Jane Smith', 'jane@example.com', 25, 'active', 150.00),
		('Bob Wilson', 'bob@example.com', 35, 'inactive', 200.00)
	`)
	assert.NoError(t, err, "Failed to insert test data")
}

// cleanupTestDB closes the database connection
func cleanupTestDB(t *testing.T, dbManager *DBManager) {
	if err := dbManager.CloseDB(); err != nil {
		t.Errorf("Failed to close database: %v", err)
	}
}

// resetTestTable truncates the test table to ensure a clean state
func resetTestTable(t *testing.T, dbManager *DBManager) {
	db, _ := dbManager.GetDB()
	_, err := db.Exec("TRUNCATE TABLE test_users")
	assert.NoError(t, err, "Failed to reset test table")
}

// testWithCleanTable is a helper function that runs a test with a clean table
func cleanTable(t *testing.T, dbManager *DBManager) {
	resetTestTable(t, dbManager)
	insertTestDataForTest(t, dbManager)
}

func TestPagination(t *testing.T) {
	dbManager := setupTestDB(t)
	defer cleanupTestDB(t, dbManager)

	qb := Table("test_users")

	expected := []map[string]interface{}{
		{"name": "John Doe", "email": "john@example.com", "age": int64(30), "price": "100.00", "status": "active"},
		{"name": "Jane Smith", "email": "jane@example.com", "age": int64(25), "price": "150.00", "status": "active"},
		{"name": "Bob Wilson", "email": "bob@example.com", "age": int64(35), "price": "200.00", "status": "inactive"},
	}

	tests := []struct {
		name     string
		data     []map[string]interface{}
		wantRows int64
		wantErr  bool
		where    map[string]interface{}
	}{
		{
			name:     "paginate all users",
			data:     expected,
			wantRows: 3,
			wantErr:  false,
		},
		{
			name: "paginate users where status is active",
			data: []map[string]interface{}{
				expected[0],
				expected[1],
			},
			wantRows: 2,
			wantErr:  false,
			where: map[string]interface{}{
				"status": "active",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			cleanTable(t, dbManager)

			qb := qb.Select("name", "email", "age", "status", "price")

			// Add the WHERE filter
			for col, val := range test.where {
				qb.Where(col, "=", val)
			}

			rows, _, err := qb.Paginate(10, 1, false)

			if test.wantErr {
				assert.Error(t, err, "Expected an error but got nil")
				return
			}
			assert.NoError(t, err, "Unexpected error")
			assert.Equal(t, test.data, rows, "Expected rows to be equal")
		})
	}

}
