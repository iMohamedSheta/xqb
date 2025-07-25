package xqb

import (
	"database/sql"
	"fmt"

	"github.com/iMohamedSheta/xqb/shared/enums"
	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// Insert inserts new rows into the database using the current connection
func (qb *QueryBuilder) Insert(values []map[string]any) (int64, error) {
	result, err := qb.insert(values)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// InsertGetId inserts new rows and returns LastInsertId using the current connection
func (qb *QueryBuilder) InsertGetId(values []map[string]any) (int64, error) {
	result, err := qb.insert(values)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (qb *QueryBuilder) Upsert(values []map[string]any, uniqueBy []string, updateColumns []string) (int64, error) {
	query, args, err := qb.UpsertSql(values, uniqueBy, updateColumns)
	if err != nil {
		return 0, err
	}

	result, err := Sql(query, args...).Connection(qb.connection).WithTx(qb.tx).Execute()
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// InsertSql returns the sql query for inserting new rows into the database
func (qb *QueryBuilder) InsertSql(values []map[string]any) (string, []any, error) {
	qb.queryType = enums.INSERT
	qbData := qb.GetData()
	qbData.InsertedValues = values

	return qb.dialect.Build(qbData)
}

// UpsertSql returns a SQL query that can be used to upsert rows into the database using the current connection
func (qb *QueryBuilder) UpsertSql(values []map[string]any, uniqueBy []string, updateColumns []string) (string, []any, error) {
	qb.queryType = enums.INSERT

	if len(values) == 0 || len(uniqueBy) == 0 || len(updateColumns) == 0 {
		return "", nil, fmt.Errorf("%w: Upsert() values, uniqueBy and updateColumns must not be empty", xqbErr.ErrInvalidQuery)
	}

	// Ensure updateColumns are part of inserted columns
	insertCols := make(map[string]struct{})
	for col := range values[0] {
		insertCols[col] = struct{}{}
	}

	for _, col := range updateColumns {
		if _, ok := insertCols[col]; !ok {
			return "", nil, fmt.Errorf("%w: Upsert() cannot update column %q because it is not part of inserted values", xqbErr.ErrInvalidQuery, col)
		}
	}

	// set upsert options
	qb.options[types.OptionIsUpsert] = true
	qb.options[types.OptionUpsertUniqueBy] = uniqueBy
	qb.options[types.OptionUpsertUpdatedCols] = updateColumns

	qbData := qb.GetData()
	qbData.InsertedValues = values

	return qb.dialect.Build(qbData)
}

// insert inserts new rows into the database using the current connection
func (qb *QueryBuilder) insert(values []map[string]any) (sql.Result, error) {
	query, args, err := qb.InsertSql(values)
	if err != nil {
		return nil, err
	}

	result, err := Sql(query, args...).Connection(qb.connection).WithTx(qb.tx).Execute()
	if err != nil {
		return nil, err
	}

	return result, nil
}
