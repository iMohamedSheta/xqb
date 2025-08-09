package xqb

import (
	"database/sql"
	"fmt"

	"github.com/iMohamedSheta/xqb/shared/enums"
	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// updates rows in the database
func (qb *QueryBuilder) Update(data map[string]any) (int64, error) {
	result, err := qb.update(data)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// core implementation of the update method
func (qb *QueryBuilder) update(data map[string]any) (sql.Result, error) {
	query, args, err := qb.UpdateSql(data)
	if err != nil {
		return nil, fmt.Errorf("%w: Update() Failed to build the sql, %v", xqbErr.ErrInvalidQuery, err)
	}

	return Sql(query, args...).
		// WithBeforeExec(qb.settings.GetOnBeforeQueryExecution()).
		WithAfterExec(qb.settings.GetOnAfterQueryExecution()).
		Connection(qb.connection).
		WithTx(qb.tx).
		Execute()
}

// UpdateSql return sql query and bindings for Update()
func (qb *QueryBuilder) UpdateSql(data map[string]any) (string, []any, error) {
	qb.queryType = enums.UPDATE

	for column, value := range data {
		binding := &types.Binding{
			Column: column,
			Value:  value,
		}
		qb.updatedBindings = append(qb.updatedBindings, binding)
	}

	return qb.ToSql()
}
