package xqb

import (
	"database/sql"
	"fmt"

	"github.com/iMohamedSheta/xqb/shared/enums"
	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// Insert inserts new rows into the database using the current connection
func (qb *QueryBuilder) Insert(values []map[string]any) error {
	_, err := qb.insert(values, false)
	if err != nil {
		return err
	}

	return nil
}

// InsertGetId inserts new rows and returns LastInsertId using the current connection
func (qb *QueryBuilder) InsertGetId(values []map[string]any) (int64, error) {
	result, err := qb.insert(values, true)
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

	result, err := Sql(query, args...).
		WithContext(qb.ctx).
		WithAfterExec(qb.settings.GetOnAfterQueryExecution()).
		Connection(qb.connection).
		WithTx(qb.tx).
		Execute()

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// InsertSql returns the sql query for inserting new rows into the database
func (qb *QueryBuilder) InsertGetIdSql(values []map[string]any) (string, []any, error) {
	return qb.insertSql(values, true)
}

// InsertSql returns the sql query for inserting new rows into the database
func (qb *QueryBuilder) InsertSql(values []map[string]any) (string, []any, error) {
	return qb.insertSql(values, false)
}

func (qb *QueryBuilder) insertSql(values []map[string]any, getId bool) (string, []any, error) {
	if getId {
		qb.SetOption(types.OptionReturningId, true)
	}

	qb.queryType = enums.INSERT
	qb.insertedValues = values
	return qb.ToSql()
}

// UpsertSql returns a Sql query that can be used to upsert rows into the database using the current connection
func (qb *QueryBuilder) UpsertSql(values []map[string]any, uniqueBy []string, updateColumns []string) (string, []any, error) {
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

	return qb.InsertSql(values)
}

// insert inserts new rows into the database using the current connection
func (qb *QueryBuilder) insert(values []map[string]any, getId bool) (sql.Result, error) {
	if getId {
		qb.SetOption(types.OptionReturningId, true)
	}

	query, args, err := qb.InsertSql(values)
	if err != nil {
		return nil, err
	}

	if getId && qb.dialect.Getdialect().String() == types.DialectPostgres.String() {
		rows, err := Sql(query, args...).
			WithContext(qb.ctx).
			WithAfterExec(qb.settings.GetOnAfterQueryExecution()).
			Connection(qb.connection).
			WithTx(qb.tx).
			Query()

		if err != nil {
			return nil, err
		}

		defer rows.Close()

		var ids []int64
		for rows.Next() {
			var id int64
			if err := rows.Scan(&id); err != nil {
				return nil, err
			}
			ids = append(ids, id)
		}

		affectedRows := int64(len(ids))

		return &QuerResult{
			lastInsertId: ids[0], // Return the first inserted ID
			affectedRows: affectedRows,
		}, nil
	}

	return Sql(query, args...).
		WithContext(qb.ctx).
		WithAfterExec(qb.settings.GetOnAfterQueryExecution()).
		Connection(qb.connection).
		WithTx(qb.tx).
		Execute()
}

type QuerResult struct {
	lastInsertId int64
	affectedRows int64
}

func (qr QuerResult) LastInsertId() (int64, error) {
	return qr.lastInsertId, nil
}

func (qr QuerResult) RowsAffected() (int64, error) {
	return qr.affectedRows, nil
}
