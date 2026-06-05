package postgres

import (
	"fmt"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileOrderByClause compiles the ORDER BY clause
func (d *PostgresDialect) compileOrderByClause(qb *types.QueryBuilderData) (string, []any, error) {
	if len(qb.OrderBy) == 0 {
		return "", nil, nil
	}

	var bindings []any
	dialect := d.Getdialect().String()
	sql := " ORDER BY "

	for i, order := range qb.OrderBy {
		if i > 0 {
			sql += ", "
		}

		if order.Raw != nil {
			// expression column: Raw or DialectExpression — dialect resolution is delegated
			// to the expression itself via ExpressionInterface.ToSql
			exprSql, exprBindings, err := order.Raw.ToSql(dialect)
			if err != nil {
				return "", nil, fmt.Errorf("%w: ORDER BY expression failed for dialect %s: %v", xqbErr.ErrInvalidQuery, dialect, err)
			}
			if exprSql == "" {
				return "", nil, fmt.Errorf("%w: ORDER BY raw SQL is empty for dialect %s", xqbErr.ErrInvalidQuery, dialect)
			}
			sql += exprSql
			bindings = append(bindings, exprBindings...)
		} else {
			// plain column name: wrap it with dialect-specific quoting
			// e.g. "created_at" → "created_at" (postgres) or `created_at` (mysql)
			sql += d.Wrap(order.Column)
		}

		// direction is empty for raw expressions like OrderByRaw(...)
		// where direction is already embedded in the SQL itself
		if order.Direction != "" {
			sql += " " + order.Direction
		}
	}

	return sql, bindings, nil
}
