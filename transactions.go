package xqb

import (
	"database/sql"
	"fmt"
)

type TransactionQuery struct {
	connection string
	tx         *sql.Tx
}

// NewTransaction - create a new transaction
func NewTransaction() *TransactionQuery {
	return &TransactionQuery{
		connection: "default",
	}
}

// Connection - set the connection
func (s *TransactionQuery) Connection(connection string) *TransactionQuery {
	if connection == "" || !Manager().HasConnection(connection) {
		connection = "default"
	}
	s.connection = connection
	return s
}

// WithTx - set the transaction
func (s *TransactionQuery) WithTx(tx *sql.Tx) *TransactionQuery {
	s.tx = tx
	return s
}

// Transaction executes a function within a transaction.
// If the function returns an error, the transaction is rolled back.
// If the function returns nil, the transaction is committed.
func Transaction(fn func(*sql.Tx) error) (err error) {
	transaction, err := BeginTransaction()
	if err != nil {
		return err
	}

	// get the transaction
	tx := transaction.tx

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
func BeginTransaction() (*TransactionQuery, error) {
	transaction := NewTransaction()
	err := transaction.BeginTransaction()
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *TransactionQuery) BeginTransaction() error {
	if t.tx != nil {
		return nil
	}

	db, err := Manager().connection(t.connection)
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	t.tx = tx
	return nil
}

// BeginTransactionWithConnection starts a new transaction and returns it.
func BeginTransactionWithConnection(connection string) (*TransactionQuery, error) {
	transaction := NewTransaction()
	transaction.Connection(connection)
	err := transaction.BeginTransaction()
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (t *TransactionQuery) Commit() error {
	if t.tx == nil {
		return nil
	}
	defer func() { t.tx = nil }()
	return t.tx.Commit()
}

// Rollback the transaction
func (t *TransactionQuery) Rollback() error {
	if t.tx == nil {
		return nil
	}
	return t.tx.Rollback()
}

// Tx returns the transaction
func (t *TransactionQuery) Tx() *sql.Tx {
	return t.tx
}
