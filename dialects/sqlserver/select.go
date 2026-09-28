package sqlserver

import (
	"fmt"
	"strings"

	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileSelectClause compiles the SELECT clause.
// When a LIMIT is set without OFFSET, SQL Server uses SELECT TOP n
// Paged limits use FETCH NEXT.
func (d *SqlServerDialect) compileSelectClause(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql string

	sql += "SELECT"

	if qb.IsUsingDistinct {
		sql += " DISTINCT"
	}

	// TOP handles limit-without-offset (OFFSET/FETCH handles the paged case)
	if qb.Limit != 0 && qb.Offset == 0 {
		sql += fmt.Sprintf(" TOP %d", qb.Limit)
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
			case types.ExpressionInterface:
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
