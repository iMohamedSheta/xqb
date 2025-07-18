package xqb

import (
	"fmt"

	"github.com/iMohamedSheta/xqb/shared/types"
)

// OrderBy adds an ORDER BY clause

func (qb *QueryBuilder) OrderBy(column any, direction string) *QueryBuilder {
	var col string
	var bindings []any

	switch v := column.(type) {
	case string:
		col = v
	case *types.Expression:
		col = v.SQL
		bindings = v.Bindings
	default:
		col = fmt.Sprintf("%v", v)
	}

	qb.orderBy = append(qb.orderBy, types.OrderBy{
		Column:    col,
		Direction: direction,
	})

	if len(bindings) > 0 {
		for _, binding := range bindings {
			qb.bindings = append(qb.bindings, types.Binding{Value: binding})
		}
	}

	return qb
}

// OrderByRaw adds a raw ORDER BY clause
func (qb *QueryBuilder) OrderByRaw(sql string, bindings ...any) *QueryBuilder {
	return qb.OrderBy(Raw(sql, bindings...), "")
}

// OrderByDesc adds an ORDER BY DESC clause
func (qb *QueryBuilder) OrderByDesc(column string) *QueryBuilder {
	return qb.OrderBy(column, "DESC")
}

// OrderByAsc adds an ORDER BY ASC clause
func (qb *QueryBuilder) OrderByAsc(column string) *QueryBuilder {
	return qb.OrderBy(column, "ASC")
}

// Latest orders by the given column in descending order
func (qb *QueryBuilder) Latest(column string) *QueryBuilder {
	return qb.OrderByDesc(column)
}

// Oldest orders by the given column in ascending order
func (qb *QueryBuilder) Oldest(column string) *QueryBuilder {
	return qb.OrderByAsc(column)
}

// InRandomOrder orders the results randomly
func (qb *QueryBuilder) InRandomOrder() *QueryBuilder {
	return qb.OrderByRaw("RAND()")
}
