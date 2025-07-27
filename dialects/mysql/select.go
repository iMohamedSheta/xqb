package mysql

import (
	"fmt"
	"strings"

	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileSelectClause compiles the SELECT clause
func (d *MySqlDialect) compileSelectClause(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql string

	sql += "SELECT"

	if qb.IsUsingDistinct {
		sql += " DISTINCT"
	}

	// Handle columns
	if len(qb.Columns) == 0 {
		sql += " *"
	} else {
		columns := make([]string, 0)

		// Add regular columns
		for _, column := range qb.Columns {
			switch v := column.(type) {
			case string:
				columns = append(columns, d.Wrap(v))
			case *types.Expression:
				columns = append(columns, v.Sql)
				bindings = append(bindings, v.Bindings...)
			case *types.DialectExpression:
				sqlStr, sqlBindings, err := v.ToSql(d.Getdialect().String())
				if err != nil {
					return "", nil, err
				}
				columns = append(columns, sqlStr)
				bindings = append(bindings, sqlBindings...)
			default:
				columns = append(columns, fmt.Sprintf("%v", v))
			}
		}

		sql += " " + strings.Join(columns, ", ")
	}

	return sql, bindings, nil
}
