package postgres

import (
	"fmt"

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
			exprSql, b, err := cte.Expression.ToSql(d.Getdialect().String())
			if err != nil {
				return "", nil, err
			}
			sql += exprSql
			bindings = append(bindings, b...)

		} else if cte.Query != nil {
			// Type assert the Query to QueryBuilderData
			// fallback: QueryBuilder support
			switch q := cte.Query.(type) {
			case *types.QueryBuilderData:
				cteSql, cteBindings, err := d.compileBaseQuery(q)
				if err != nil {
					return "", nil, err
				}
				sql += cteSql
				bindings = append(bindings, cteBindings...)

			default:
				return "", nil, fmt.Errorf("unsupported CTE query type: %T", cte.Query)
			}
		}

		sql += ")"
	}

	return sql + " ", bindings, nil
}
