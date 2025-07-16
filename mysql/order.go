package mysql

import (
	"strings"

	"github.com/iMohamedSheta/xqb/types"
)

// compileOrderByClause compiles the ORDER BY clause
func (mg *MySQLGrammar) compileOrderByClause(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql strings.Builder

	if len(qb.OrderBy) > 0 {
		sql.WriteString(" ORDER BY ")
		for i, order := range qb.OrderBy {
			if i > 0 {
				sql.WriteString(", ")
			}
			sql.WriteString(order.Column)
			sql.WriteString(" ")
			sql.WriteString(order.Direction)
		}
	}

	return sql.String(), bindings, nil
}
