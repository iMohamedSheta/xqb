package xqb

import (
	"database/sql"
	"fmt"

	"github.com/iMohamedSheta/xqb/shared/enums"
)

// insert inserts new rows into the database using the current connection
func (qb *QueryBuilder) insert(data []map[string]any) (sql.Result, error) {
	db, err := Connection(qb.connection)
	if err != nil {
		return nil, err
	}

	qb.queryType = enums.INSERT
	qbData := qb.GetData()
	qbData.InsertedValues = data

	query, args, err := qb.grammar.Build(qbData)
	if err != nil {
		return nil, err
	}

	var result sql.Result
	if qb.tx != nil {
		result, err = qb.tx.Exec(query, args...)
	} else {
		result, err = db.Exec(query, args...)
	}

	if err != nil {
		return nil, fmt.Errorf("insert failed: %w", err)
	}

	return result, nil
}

// Insert inserts new rows into the database using the current connection
func (qb *QueryBuilder) Insert(data []map[string]any) (int64, error) {
	result, err := qb.insert(data)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// InsertGetId inserts new rows and returns LastInsertId using the current connection
func (qb *QueryBuilder) InsertGetId(data []map[string]any) (int64, error) {
	result, err := qb.insert(data)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}
