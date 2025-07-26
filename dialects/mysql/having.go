package mysql

import (
	"strings"

	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileHavingClause compiles the HAVING clause
func (mg *MySqlDialect) compileHavingClause(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql strings.Builder

	if len(qb.Having) > 0 {
		sql.WriteString(" HAVING ")

		for i, having := range qb.Having {
			if i > 0 {
				sql.WriteString(" ")
				sql.WriteString(string(having.Connector))
				sql.WriteString(" ")
			}

			// Use raw expression if available
			if having.Raw != nil {
				sql.WriteString(having.Raw.Sql)
				bindings = append(bindings, having.Raw.Bindings...)
			} else {
				sql.WriteString(mg.Wrap(having.Column))
				sql.WriteString(" ")
				sql.WriteString(having.Operator)
				if having.Value != nil {
					sql.WriteString(" ?")
					bindings = append(bindings, having.Value)
				}
			}
		}
	}

	return sql.String(), bindings, nil
}
