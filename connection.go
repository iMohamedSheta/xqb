package xqb

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrNoConnection = errors.New("no database connection established")
)

// DBManager handles database connection management
type DBManager struct {
	db *sql.DB
	mu sync.RWMutex
}

var (
	dbManagerInstance *DBManager
	dbManagerOnce     sync.Once
)

// GetDBManager returns the singleton instance of DBManager
func GetDBManager() *DBManager {
	dbManagerOnce.Do(func() {
		dbManagerInstance = &DBManager{}
	})
	return dbManagerInstance
}

// SetDB sets the database connection (sql.DB)
func (m *DBManager) SetDB(db *sql.DB) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.db = db
}

// GetDB returns the current database connection (sql.DB)
func (m *DBManager) GetDB() (*sql.DB, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.db == nil {
		return nil, ErrNoConnection
	}
	return m.db, nil
}

// Transaction executes a function within a transaction.
// If the function returns an error, the transaction is rolled back.
// If the function returns nil, the transaction is committed.
func (m *DBManager) Transaction(fn func(*sql.Tx) error) (err error) {
	db, err := m.GetDB()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Convert panic to error
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			switch e := p.(type) {
			case error:
				err = fmt.Errorf("transaction error: %w", e)
			default:
				err = fmt.Errorf("transaction error: %v", p)
			}
		}
	}()

	// Execute the function
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

// BeginTransaction starts a new transaction and returns it.
func (m *DBManager) BeginTransaction() (*sql.Tx, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	return db.Begin()
}

// CloseDB closes the database connection
func (m *DBManager) CloseDB() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.db == nil {
		return nil
	}

	err := m.db.Close()
	m.db = nil
	return err
}

// IsDBConnected checks if there is an active database connection
func (m *DBManager) IsDBConnected() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.db != nil
}
