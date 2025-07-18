package xqb

import (
	"fmt"

	"github.com/iMohamedSheta/xqb/shared/types"
)

// GroupBy sets the GROUP BY clause
func (qb *QueryBuilder) GroupBy(columns ...any) *QueryBuilder {
	for _, column := range columns {
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

		qb.groupBy = append(qb.groupBy, col)

		if len(bindings) > 0 {
			for _, binding := range bindings {
				qb.bindings = append(qb.bindings, types.Binding{Value: binding})
			}
		}
	}
	return qb
}

// GroupByRaw adds a raw GROUP BY clause
func (qb *QueryBuilder) GroupByRaw(sql string, bindings ...any) *QueryBuilder {
	return qb.GroupBy(Raw(sql, bindings...))
}

// Having adds a HAVING clause
func (qb *QueryBuilder) Having(column any, operator string, value any) *QueryBuilder {
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

	qb.having = append(qb.having, types.Having{
		Column:    col,
		Operator:  operator,
		Value:     value,
		Connector: types.AND,
	})

	if len(bindings) > 0 {
		for _, binding := range bindings {
			qb.bindings = append(qb.bindings, types.Binding{Value: binding})
		}
	}

	return qb
}

// HavingRaw adds a raw HAVING clause
func (qb *QueryBuilder) HavingRaw(sql string, bindings ...any) *QueryBuilder {
	return qb.Having(Raw(sql, bindings...), "", nil)
}

// OrHaving adds an OR HAVING clause
func (qb *QueryBuilder) OrHaving(column any, operator string, value any) *QueryBuilder {
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

	qb.having = append(qb.having, types.Having{
		Column:    col,
		Operator:  operator,
		Value:     value,
		Connector: types.OR,
	})

	if len(bindings) > 0 {
		for _, binding := range bindings {
			qb.bindings = append(qb.bindings, types.Binding{Value: binding})
		}
	}

	return qb
}

// OrHavingRaw adds a raw OR HAVING clause
func (qb *QueryBuilder) OrHavingRaw(sql string, bindings ...any) *QueryBuilder {
	return qb.OrHaving(Raw(sql, bindings...), "", nil)
}
