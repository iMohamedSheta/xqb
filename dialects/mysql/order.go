package mysql

import (
	"fmt"
	"strings"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileOrderByClause compiles the ORDER BY clause
func (mg *MySqlDialect) compileOrderByClause(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql strings.Builder

	if len(qb.OrderBy) > 0 {
		sql.WriteString(" ORDER BY ")
		for i, order := range qb.OrderBy {
			if i > 0 {
				sql.WriteString(", ")
			}
			if order.Raw != nil {
				expr := order.Raw.Dialects[mg.GetDriver().String()]
				if expr == nil {
					expr = order.Raw.Dialects[order.Raw.Default]
				}

				if expr == nil {
					return "", nil, fmt.Errorf("%w: ORDER BY raw Sql not supported for %s dialect you need to specify ORDER BY column the dialectExpression", xqbErr.ErrInvalidQuery, mg.GetDriver().String())
				}

				sql.WriteString(expr.Sql)
				bindings = append(bindings, expr.Bindings...)
			} else {
				sql.WriteString(mg.Wrap(order.Column))
			}

			if order.Direction != "" {
				sql.WriteString(" ")
				sql.WriteString(order.Direction)
			}
		}
	}

	return sql.String(), bindings, nil
}
