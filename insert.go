package xqb

import (
	"database/sql"
	"fmt"

	"github.com/iMohamedSheta/xqb/types"
)

// Insert inserts new rows into the database using the current connection
func (qb *QueryBuilder) Insert(data []map[string]interface{}) (int64, error) {
	dbManager := GetDBManager()
	if !dbManager.IsDBConnected() {
		return 0, ErrNoConnection
	}

	db, err := dbManager.GetDB()
	if err != nil {
		return 0, err
	}

	qb.queryType = types.INSERT
	qbData := qb.GetData()
	qbData.InsertedValues = data

	query, args, err := qb.grammar.Build(qbData)
	if err != nil {
		return 0, err
	}

	result, err := db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("insert failed: %w", err)
	}

	return result.RowsAffected()
}

// InsertWithTx inserts new rows into the database using the provided transaction
func (qb *QueryBuilder) InsertWithTx(data []map[string]interface{}, tx *sql.Tx) (int64, error) {
	if tx == nil {
		return 0, fmt.Errorf("transaction is required")
	}

	qb.queryType = types.INSERT
	qbData := qb.GetData()
	qbData.InsertedValues = data

	query, args, err := qb.grammar.Build(qbData)
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("insert failed: %w", err)
	}

	return result.RowsAffected()
}

// InsertGetID inserts new rows and returns LastInsertedID using the current connection
func (qb *QueryBuilder) InsertGetId(data []map[string]interface{}) (int64, error) {
	dbManager := GetDBManager()
	if !dbManager.IsDBConnected() {
		return 0, ErrNoConnection
	}

	db, err := dbManager.GetDB()
	if err != nil {
		return 0, err
	}

	qb.queryType = types.INSERT
	qbData := qb.GetData()
	qbData.InsertedValues = data

	query, args, err := qb.grammar.Build(qbData)
	if err != nil {
		return 0, err
	}

	result, err := db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("insert failed: %w", err)
	}

	return result.LastInsertId()
}

// InsertGetIdWithTx inserts new rows and returns LastInsertedID using the provided transaction
func (qb *QueryBuilder) InsertGetIdWithTx(data []map[string]interface{}, tx *sql.Tx) (int64, error) {
	if tx == nil {
		return 0, fmt.Errorf("transaction is required")
	}

	qb.queryType = types.INSERT
	qbData := qb.GetData()
	qbData.InsertedValues = data

	query, args, err := qb.grammar.Build(qbData)
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("insert failed: %w", err)
	}

	return result.LastInsertId()
}
