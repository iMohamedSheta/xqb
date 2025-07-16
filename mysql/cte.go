package mysql

import (
	"strings"

	"github.com/iMohamedSheta/xqb/types"
)

// compileCTEs compiles Common Table Expressions
func (mg *MySQLGrammar) compileCTEs(qb *types.QueryBuilderData) (string, []any, error) {
	if len(qb.WithCTEs) == 0 {
		return "", nil, nil
	}

	var bindings []any
	var sql strings.Builder

	sql.WriteString("WITH ")
	for i, cte := range qb.WithCTEs {
		if i > 0 {
			sql.WriteString(", ")
		}
		sql.WriteString(cte.Name)
		sql.WriteString(" AS (")

		if cte.Expression != nil {
			// Use raw expression if provided
			sql.WriteString(cte.Expression.SQL)
			bindings = append(bindings, cte.Expression.Bindings...)
		} else if cte.Query != nil {
			// Type assert the Query to QueryBuilderData
			if queryData, ok := cte.Query.(*types.QueryBuilderData); ok {
				cteSQL, cteBindings, _ := mg.compileBaseQuery(queryData)
				sql.WriteString(cteSQL)
				bindings = append(bindings, cteBindings...)
			}
		}

		sql.WriteString(")")
	}

	return sql.String(), bindings, nil
}
