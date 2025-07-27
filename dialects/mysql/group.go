package mysql

import (
	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileGroupByClause compiles the GROUP BY clause
func (d *MySqlDialect) compileGroupByClause(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql string

	if len(qb.GroupBy) > 0 {
		sql += " GROUP BY "
		for i, column := range qb.GroupBy {
			if i > 0 {
				sql += ", "
			}
			sql += d.Wrap(column)
		}
	}

	return sql, bindings, nil
}
