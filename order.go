package xqb

import (
	"fmt"

	"github.com/iMohamedSheta/xqb/shared/types"
)

// OrderBy adds an ORDER BY clause
func (qb *QueryBuilder) orderByClause(column any, direction string) *QueryBuilder {
	var col string
	var raw types.ExpressionInterface

	switch v := column.(type) {
	case string:
		// plain column name: OrderBy("created_at", "DESC")
		col = v
	case types.ExpressionInterface:
		// expression as column: OrderBy(Raw("FIELD(status, 'active', 'inactive')"), "ASC")
		// both *types.Expression and *types.DialectExpression are covered here
		// since both implement ExpressionInterface — no wrapping needed
		raw = v
	default:
		// fallback: stringify anything unexpected
		col = fmt.Sprintf("%v", v)
	}

	qb.orderBy = append(qb.orderBy, &types.OrderBy{
		Column:    col,
		Direction: direction,
		Raw:       raw,
	})

	return qb
}

// OrderBy adds an ORDER BY clause with an explicit direction.
// column can be a plain string or any types.ExpressionInterface (Raw, DialectExpression, etc.)
func (qb *QueryBuilder) OrderBy(column any, direction string) *QueryBuilder {
	return qb.orderByClause(column, direction)
}

// OrderByRaw adds a raw ORDER BY clause
func (qb *QueryBuilder) OrderByRaw(sql string, bindings ...any) *QueryBuilder {
	return qb.orderByClause(Raw(sql, bindings...), "")
}

// OrderByDesc adds an ORDER BY DESC clause
func (qb *QueryBuilder) OrderByDesc(column string) *QueryBuilder {
	return qb.orderByClause(column, "DESC")
}

// OrderByAsc adds an ORDER BY ASC clause
func (qb *QueryBuilder) OrderByAsc(column string) *QueryBuilder {
	return qb.orderByClause(column, "ASC")
}

// Latest orders by the given column in descending order
func (qb *QueryBuilder) Latest(column string) *QueryBuilder {
	return qb.OrderByDesc(column)
}

// Oldest orders by the given column in ascending order
func (qb *QueryBuilder) Oldest(column string) *QueryBuilder {
	return qb.OrderByAsc(column)
}
