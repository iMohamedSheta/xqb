package mysql

import (
	"fmt"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileOrderByClause compiles the ORDER BY clause
func (d *MySqlDialect) compileOrderByClause(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql string

	if len(qb.OrderBy) > 0 {
		sql += " ORDER BY "
		for i, order := range qb.OrderBy {
			if i > 0 {
				sql += ", "
			}
			if order.Raw != nil {
				expr := order.Raw.Dialects[d.Getdialect().String()]
				if expr == nil {
					expr = order.Raw.Dialects[order.Raw.Default]
				}

				if expr == nil {
					return "", nil, fmt.Errorf("%w: ORDER BY raw Sql not supported for %s dialect you need to specify ORDER BY column the dialectExpression", xqbErr.ErrInvalidQuery, d.Getdialect().String())
				}

				sql += expr.Sql
				bindings = append(bindings, expr.Bindings...)
			} else {
				sql += d.Wrap(order.Column)
			}

			if order.Direction != "" {
				sql += " " + order.Direction
			}
		}
	}

	return sql, bindings, nil
}
