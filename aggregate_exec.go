package xqb

import (
	"database/sql"
	"fmt"
)

// Count - Returns the number of rows in the result set.
func (qb *QueryBuilder) Count(column string, tx *sql.Tx) (int64, error) {
	qb.columns = []any{fmt.Sprintf("COUNT(%s) as count_value", column)}

	data, err := qb.Execute(tx)
	if err != nil {
		return 0, err
	}

	if len(data) != 1 {
		return 0, fmt.Errorf("expected 1 row, got %d", len(data))
	}

	count, ok := data[0]["count_value"].(int64)
	if !ok {
		return 0, fmt.Errorf("expected count to be uint64, got %T", data[0]["count_value"])
	}

	return count, nil
}
