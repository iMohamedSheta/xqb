package postgres

import (
	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileCTEs compiles Common Table Expressions
func (d *PostgresDialect) compileCTEs(qb *types.QueryBuilderData) (string, []any, error) {
	if len(qb.WithCTEs) == 0 {
		return "", nil, nil
	}

	var bindings []any
	var sql string

	sql += "WITH "
	for i, cte := range qb.WithCTEs {
		if i > 0 {
			sql += ", "
		}
		if cte.Recursive {
			sql += "RECURSIVE "
		}
		sql += cte.Name + " AS ("

		if cte.Expression != nil {
			// Use raw expression if provided
			sql += cte.Expression.Sql
			bindings = append(bindings, cte.Expression.Bindings...)
		} else if cte.Query != nil {
			// Type assert the Query to QueryBuilderData
			if queryData, ok := cte.Query.(*types.QueryBuilderData); ok {
				cteSql, cteBindings, err := d.compileBaseQuery(queryData)
				if err != nil {
					return "", nil, err
				}
				sql += cteSql
				bindings = append(bindings, cteBindings...)
			}
		}

		sql += ")"
	}

	return sql + " ", bindings, nil
}
