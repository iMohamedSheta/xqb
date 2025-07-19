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

// core delete execution method
func (qb *QueryBuilder) delete() (int64, error) {
	qb.queryType = enums.DELETE
	qbData := qb.GetData()

	query, args, err := qb.grammar.Build(qbData)
	if err != nil {
		return 0, err
	}

	var result sql.Result

	if qb.tx != nil {
		result, err = qb.tx.Exec(query, args...)
		if err != nil {
			return 0, fmt.Errorf("delete failed:  %w", err)
		}

	} else {
		db, err := Connection(qb.connection)
		if err != nil {
			return 0, err
		}

		result, err = db.Exec(query, args...)
		if err != nil {
			return 0, fmt.Errorf("delete failed: %w", err)
		}
	}

	return result.RowsAffected()
}
