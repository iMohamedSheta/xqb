package postgres

import (
	"strings"

	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileGroupByClause compiles the GROUP BY clause
func (pg *PostgresGrammar) compileGroupByClause(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql strings.Builder

	if len(qb.GroupBy) > 0 {
		sql.WriteString(" GROUP BY ")
		for i, column := range qb.GroupBy {
			if i > 0 {
				sql.WriteString(", ")
			}
			sql.WriteString(column)
		}
	}

	return sql.String(), bindings, nil
}
