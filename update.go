package xqb

import (
	"database/sql"
	"fmt"

	"github.com/iMohamedSheta/xqb/types"
)

// updates rows in the database
func (qb *QueryBuilder) Update(data map[string]any) (int64, error) {
	return qb.update(data, nil)
}

// InsertTx inserts new rows into the database using the provided transaction
func (qb *QueryBuilder) UpdateTx(data map[string]interface{}, tx *sql.Tx) (int64, error) {
	return qb.update(data, tx)
}

// core implementation of the update method
func (qb *QueryBuilder) update(data map[string]interface{}, tx *sql.Tx) (int64, error) {

	qb.queryType = types.UPDATE
	qbData := qb.GetData()

	for column, value := range data {
		binding := types.Binding{
			Column: column,
			Value:  value,
		}
		qb.bindings = append(qb.bindings, binding)
	}

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
