package xqb

import (
	"database/sql"
	"fmt"

	"github.com/iMohamedSheta/xqb/types"
)

// Delete rows in the database
func (qb *QueryBuilder) Delete() (int64, error) {
	return qb.delete(nil)
}

// Delete rows in the database with transaction
func (qb *QueryBuilder) DeleteTx(tx *sql.Tx) (int64, error) {
	return qb.delete(tx)
}

// core delete execution method
func (qb *QueryBuilder) delete(tx *sql.Tx) (int64, error) {
	qb.queryType = types.DELETE
	qbData := qb.GetData()

	query, args, err := qb.grammar.Build(qbData)
	if err != nil {
		return 0, err
	}

	var result sql.Result

	if tx != nil {
		result, err = tx.Exec(query, args...)
		if err != nil {
			return 0, fmt.Errorf("update failed:  %w", err)
		}

	} else {
		dbManager := GetDBManager()
		if !dbManager.IsDBConnected() {
			return 0, ErrNoConnection
		}

		db, err := dbManager.GetDB()

		if err != nil {
			return 0, err
		}

		result, err = db.Exec(query, args...)

		if err != nil {
			return 0, fmt.Errorf("update failed: %w", err)
		}
	}

	return result.RowsAffected()
}
