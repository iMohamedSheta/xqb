package xqb

import (
	"database/sql"
	"fmt"

	"github.com/iMohamedSheta/xqb/shared/enums"
)

// Delete rows in the database
func (qb *QueryBuilder) Delete() (int64, error) {
	return qb.delete()
}

// DeleteSql returns the sql query for deleting rows
func (qb *QueryBuilder) DeleteSql(table ...string) (string, []any, error) {
	if len(table) != 0 {
		qb.deleteFrom = table
	}

	qb.queryType = enums.DELETE

	return qb.ToSql()
}

// core delete execution method
func (qb *QueryBuilder) delete(table ...string) (int64, error) {
	query, args, err := qb.DeleteSql(table...)
	if err != nil {
		return 0, err
	}

	var result sql.Result

	result, err = Sql(query, args...).
		WithContext(qb.ctx).
		WithAfterExec(qb.settings.GetOnAfterQueryExecution()).
		Connection(qb.connection).
		WithTx(qb.tx).
		Execute()

	if err != nil {
		return 0, fmt.Errorf("delete failed:  %w", err)
	}

	return result.RowsAffected()
}
