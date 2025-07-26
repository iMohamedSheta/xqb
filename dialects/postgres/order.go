package postgres

import (
	"fmt"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileOrderByClause compiles the ORDER BY clause
func (pg *PostgresDialect) compileOrderByClause(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql string

	if len(qb.OrderBy) > 0 {
		sql += " ORDER BY "
		for i, order := range qb.OrderBy {
			if i > 0 {
				sql += ", "
			}
			if order.Raw != nil {
				expr := order.Raw.Dialects[pg.GetDriver().String()]
				if expr == nil {
					expr = order.Raw.Dialects[order.Raw.Default]
				}

				if expr == nil {
					return "", nil, fmt.Errorf("%w: ORDER BY raw Sql not supported for %s dialect you need to specify ORDER BY column the dialectExpression", xqbErr.ErrInvalidQuery, pg.GetDriver().String())
				}

				sql += expr.Sql
				bindings = append(bindings, expr.Bindings...)
			} else {
				sql += pg.Wrap(order.Column)
			}

			if order.Direction != "" {
				sql += " " + order.Direction
			}
		}
	}

	return sql, bindings, nil
}
