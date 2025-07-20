package xqb

import (
	"database/sql"
	"fmt"
)

// BeginTx starts a transaction using the default connection.
func BeginTx() (*sql.Tx, error) {
	return BeginTxOn("default")
}

// BeginTxOn starts a transaction using the specified connection.
func BeginTxOn(connection string) (*sql.Tx, error) {
	if !DBManager().HasConnection(connection) {
		return nil, fmt.Errorf("%w: invalid connection %s", ErrNoConnection, connection)
	}

	db, err := DBManager().Connection(connection)
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// Transaction runs a function inside a transaction on the default connection.
func Transaction(fn func(*sql.Tx) error) error {
	return TransactionOn("default", fn)
}

// TransactionOn runs a function inside a transaction on the given connection.
func TransactionOn(connection string, fn func(*sql.Tx) error) (err error) {
	tx, err := BeginTxOn(connection)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			switch e := p.(type) {
			case error:
				err = fmt.Errorf("%w: %v", ErrTransactionFailed, e)
			default:
				err = fmt.Errorf("%w: panic %v", ErrTransactionFailed, p)
			}
		}
	}()

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
