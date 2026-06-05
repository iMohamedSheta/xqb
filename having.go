package xqb

import (
	"fmt"

	"github.com/iMohamedSheta/xqb/shared/types"
)

// Having adds a HAVING clause
func (qb *QueryBuilder) Having(column any, operator string, value any) *QueryBuilder {
	return qb.havingClause(column, operator, value, types.AND)
}

// HavingRaw adds a raw HAVING clause
func (qb *QueryBuilder) HavingRaw(sql string, bindings ...any) *QueryBuilder {
	return qb.havingClause(Raw(sql, bindings...), "", nil, types.AND)
}

// OrHaving adds an OR HAVING clause
func (qb *QueryBuilder) OrHaving(column any, operator string, value any) *QueryBuilder {
	return qb.havingClause(column, operator, value, types.OR)
}

// OrHavingRaw adds a raw OR HAVING clause
func (qb *QueryBuilder) OrHavingRaw(sql string, bindings ...any) *QueryBuilder {
	return qb.havingClause(Raw(sql, bindings...), "", nil, types.OR)
}
func (qb *QueryBuilder) havingClause(column any, operator string, value any, connector types.WhereConditionEnum) *QueryBuilder {
	var col string
	var raw types.ExpressionInterface
	var extraBindings []any

	switch v := column.(type) {
	case string:
		col = v
	case types.ExpressionInterface:
		exprSQL, exprBindings, err := v.ToSql(qb.GetDialect().Getdialect().String())
		if err != nil {
			qb.appendError(err)
			return qb
		}

		sql := exprSQL
		if operator != "" {
			sql = fmt.Sprintf("%s %s", exprSQL, operator)
			switch val := value.(type) {
			case types.ExpressionInterface:
				valSQL, valBindings, err := val.ToSql(qb.GetDialect().Getdialect().String())
				if err != nil {
					qb.appendError(err)
					return qb
				}
				sql = fmt.Sprintf("%s %s", sql, valSQL)
				extraBindings = append(extraBindings, valBindings...)
			default:
				if val != nil {
					sql = fmt.Sprintf("%s ?", sql)
					extraBindings = append(extraBindings, val)
				}
			}
		}

		raw = &types.Expression{
			Sql:      sql,
			Bindings: append(exprBindings, extraBindings...),
		}
	default:
		col = fmt.Sprintf("%v", v)
	}

	qb.having = append(qb.having, &types.Having{
		Column:    col,
		Operator:  operator,
		Value:     value,
		Connector: connector,
		Raw:       raw,
	})

	return qb
}
