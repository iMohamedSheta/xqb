package xqb

import (
	"database/sql"
	"fmt"

	"github.com/iMohamedSheta/xqb/shared/enums"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// updates rows in the database
func (qb *QueryBuilder) Update(data map[string]any) (int64, error) {
	return qb.update(data)
}

// core implementation of the update method
func (qb *QueryBuilder) update(data map[string]any) (int64, error) {

	qb.queryType = enums.UPDATE
	qbData := qb.GetData()

	for column, value := range data {
		binding := types.Binding{
			Column: column,
			Value:  value,
		}
		qbData.UpdatedBindings = append(qbData.UpdatedBindings, binding)
	}

	query, args, err := qb.grammar.Build(qbData)
	if err != nil {
		return 0, err
	}

	var result sql.Result

	if qb.tx != nil {
		result, err = qb.tx.Exec(query, args...)
		if err != nil {
			return 0, fmt.Errorf("update failed:  %w", err)
		}

	} else {
		db, err := Connection(qb.connection)
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
