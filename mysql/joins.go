package mysql

import (
	"strings"

	"github.com/iMohamedSheta/xqb/types"
)

// compileJoins compiles the JOIN clause
func (mg *MySQLGrammar) compileJoins(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql strings.Builder

	if len(qb.Joins) > 0 {
		for _, join := range qb.Joins {
			sql.WriteString(" ")
			sql.WriteString(string(join.Type))
			sql.WriteString(" ")
			sql.WriteString(join.Table)
			sql.WriteString(" ON ")
			sql.WriteString(join.Condition)
		}
	}

	return sql.String(), bindings, nil
}
