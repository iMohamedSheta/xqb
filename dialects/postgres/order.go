package postgres

import (
	"strings"

	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileOrderByClause compiles the ORDER BY clause
func (pg *PostgresDialect) compileOrderByClause(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql strings.Builder

	if len(qb.OrderBy) > 0 {
		sql.WriteString(" ORDER BY ")
		for i, order := range qb.OrderBy {
			if i > 0 {
				sql.WriteString(", ")
			}
			if order.Raw != nil {
				expr := order.Raw.Dialects[pg.GetDriver().String()]

				sql.WriteString(expr.SQL)
				bindings = append(bindings, expr.Bindings...)
			} else {
				sql.WriteString(order.Column)
			}

			if order.Direction != "" {
				sql.WriteString(" ")
				sql.WriteString(order.Direction)
			}
		}
	}

	return sql.String(), bindings, nil
}
