package xqb

import (
	"database/sql"

	"github.com/iMohamedSheta/xqb/shared/enums"
)

// insert inserts new rows into the database using the current connection
func (qb *QueryBuilder) insert(data []map[string]any) (sql.Result, error) {
	qb.queryType = enums.INSERT
	qbData := qb.GetData()
	qbData.InsertedValues = data

	query, args, err := qb.dialect.Build(qbData)
	if err != nil {
		return nil, err
	}

	result, err := Sql(query, args...).Connection(qb.connection).WithTx(qb.tx).Execute()
	if err != nil {
		return nil, err
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
