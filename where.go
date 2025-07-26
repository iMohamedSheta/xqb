package xqb

import (
	"fmt"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

func (qb *QueryBuilder) whereClause(column any, operator string, value any, connector types.WhereConditionEnum) *QueryBuilder {
	var col string
	var raw *types.Expression
	var bindings []any

	switch v := column.(type) {
	case string:
		col = v
	case *types.Expression:
		switch val := value.(type) {
		case *types.Expression:
			raw = &types.Expression{
				Sql:      fmt.Sprintf("(%s) %s (%s)", v.Sql, operator, val.Sql),
				Bindings: append(v.Bindings, val.Bindings...),
			}
		default:
			raw = &types.Expression{
				Sql:      fmt.Sprintf("%s %s ?", v.Sql, operator),
				Bindings: append(v.Bindings, val),
			}
		}
	default:
		col = fmt.Sprintf("%v", v)
	}

	// Subquery or expression as value
	if raw == nil {
		switch v := value.(type) {
		case *QueryBuilder:
			subSql, subBindings, err := v.SetDialect(qb.GetDialect().GetDriver()).ToSql()
			if err != nil {
				qb.appendError(err)
			}
			raw = &types.Expression{
				Sql:      fmt.Sprintf("%s %s (%s)", col, operator, subSql),
				Bindings: subBindings,
			}
		case *types.Expression:
			raw = &types.Expression{
				Sql:      fmt.Sprintf("%s %s (%s)", col, operator, v.Sql),
				Bindings: v.Bindings,
			}
		default:
			bindings = append(bindings, v)
		}
	}

	// Add bindings
	for _, b := range bindings {
		qb.bindings = append(qb.bindings, &types.Binding{Value: b})
	}
	if raw != nil {
		qb.bindings = append(qb.bindings, toBindings(raw.Bindings)...)
	}

	qb.where = append(qb.where, &types.WhereCondition{
		Column:    col,
		Operator:  operator,
		Value:     value,
		Connector: connector,
		Raw:       raw,
	})
	return qb
}

func (qb *QueryBuilder) Where(column any, operator string, value any) *QueryBuilder {
	return qb.whereClause(column, operator, value, types.AND)
}

// OrWhere adds an OR WHERE condition
func (qb *QueryBuilder) OrWhere(column any, operator string, value any) *QueryBuilder {
	return qb.whereClause(column, operator, value, types.OR)
}

func toBindings(vals []any) []*types.Binding {
	bs := make([]*types.Binding, len(vals))
	for i, v := range vals {
		bs[i] = &types.Binding{Value: v}
	}
	return bs
}

func (qb *QueryBuilder) whereColumnClause(column, operator string, value any, connector types.WhereConditionEnum) *QueryBuilder {
	qb.where = append(qb.where, &types.WhereCondition{
		Column:    column,
		Operator:  operator,
		Value:     value,
		Connector: connector,
	})

	return qb
}

func (qb *QueryBuilder) WhereValue(column string, operator string, value any) *QueryBuilder {
	return qb.whereColumnClause(column, operator, value, types.AND)
}

func (qb *QueryBuilder) OrWhereValue(column string, operator string, value any) *QueryBuilder {
	return qb.whereColumnClause(column, operator, value, types.OR)
}

func (qb *QueryBuilder) WhereSub(column string, operator string, sub *QueryBuilder) *QueryBuilder {
	subSql, subBindings, _ := sub.SetDialect(qb.GetDialect().GetDriver()).ToSql()
	qb.where = append(qb.where, &types.WhereCondition{
		Column:    column,
		Operator:  operator,
		Connector: types.AND,
		Raw: &types.Expression{
			Sql:      fmt.Sprintf("%s %s (%s)", column, operator, subSql),
			Bindings: subBindings,
		},
	})
	return qb
}

func (qb *QueryBuilder) OrWhereSub(column string, operator string, sub *QueryBuilder) *QueryBuilder {
	subSql, subBindings, _ := sub.SetDialect(qb.GetDialect().GetDriver()).ToSql()
	qb.where = append(qb.where, &types.WhereCondition{
		Column:    column,
		Operator:  operator,
		Connector: types.OR,
		Raw: &types.Expression{
			Sql:      fmt.Sprintf("%s %s (%s)", column, operator, subSql),
			Bindings: subBindings,
		},
	})
	return qb
}

func (qb *QueryBuilder) WhereExpr(column string, operator string, expr *types.Expression) *QueryBuilder {
	qb.where = append(qb.where, &types.WhereCondition{
		Column:    column,
		Operator:  operator,
		Connector: types.AND,
		Raw: &types.Expression{
			Sql:      fmt.Sprintf("%s %s (%s)", column, operator, expr.Sql),
			Bindings: expr.Bindings,
		},
	})
	return qb
}

func (qb *QueryBuilder) OrWhereExpr(column string, operator string, expr *types.Expression) *QueryBuilder {
	qb.where = append(qb.where, &types.WhereCondition{
		Column:    column,
		Operator:  operator,
		Connector: types.OR,
		Raw: &types.Expression{
			Sql:      fmt.Sprintf("%s %s (%s)", column, operator, expr.Sql),
			Bindings: expr.Bindings,
		},
	})
	return qb
}

func (qb *QueryBuilder) whereRawClause(sql string, bindings []any, connector types.WhereConditionEnum) *QueryBuilder {
	qb.where = append(qb.where, &types.WhereCondition{
		Raw: &types.Expression{
			Sql:      sql,
			Bindings: bindings,
		},
		Connector: connector,
	})
	qb.bindings = append(qb.bindings, toBindings(bindings)...)
	return qb
}

// WhereRaw adds a raw WHERE condition
func (qb *QueryBuilder) WhereRaw(sql string, bindings ...any) *QueryBuilder {
	return qb.whereRawClause(sql, bindings, types.AND)
}

// OrWhereRaw adds a raw OR WHERE condition
func (qb *QueryBuilder) OrWhereRaw(sql string, bindings ...any) *QueryBuilder {
	return qb.whereRawClause(sql, bindings, types.OR)
}

// WhereNull adds a WHERE NULL condition
func (qb *QueryBuilder) WhereNull(column string) *QueryBuilder {
	qb.where = append(qb.where, &types.WhereCondition{
		Column:    column,
		Operator:  "IS NULL",
		Value:     nil,
		Connector: types.AND,
	})
	return qb
}

// OrWhereNull adds an OR WHERE NULL condition
func (qb *QueryBuilder) OrWhereNull(column string) *QueryBuilder {
	qb.where = append(qb.where, &types.WhereCondition{
		Column:    column,
		Operator:  "IS NULL",
		Value:     nil,
		Connector: types.OR,
	})
	return qb
}

// WhereNotNull adds a WHERE NOT NULL condition
func (qb *QueryBuilder) WhereNotNull(column string) *QueryBuilder {
	qb.where = append(qb.where, &types.WhereCondition{
		Column:    column,
		Operator:  "IS NOT NULL",
		Value:     nil,
		Connector: types.AND,
	})
	return qb
}

// OrWhereNotNull adds an OR WHERE NOT NULL condition
func (qb *QueryBuilder) OrWhereNotNull(column string) *QueryBuilder {
	qb.where = append(qb.where, &types.WhereCondition{
		Column:    column,
		Operator:  "IS NOT NULL",
		Value:     nil,
		Connector: types.OR,
	})
	return qb
}

func (qb *QueryBuilder) whereInClause(column string, values []any, operator string, connector types.WhereConditionEnum) *QueryBuilder {
	if len(values) == 0 {
		return qb
	}

	// If single value is a subquery or expression
	for _, value := range values {
		switch v := value.(type) {
		case *QueryBuilder:
			subSql, subBindings, err := v.SetDialect(qb.GetDialect().GetDriver()).ToSql()
			if err != nil {
				qb.appendError(err)
			}
			qb.where = append(qb.where, &types.WhereCondition{
				Column:    column,
				Operator:  operator,
				Value:     nil,
				Connector: connector,
				Raw: &types.Expression{
					Sql:      fmt.Sprintf("%s %s (%s)", column, operator, subSql),
					Bindings: subBindings,
				},
			})
			return qb

		case *types.Expression:
			qb.where = append(qb.where, &types.WhereCondition{
				Column:    column,
				Operator:  operator,
				Value:     nil,
				Connector: connector,
				Raw: &types.Expression{
					Sql:      fmt.Sprintf("%s %s (%s)", column, operator, v.Sql),
					Bindings: v.Bindings,
				},
			})
			return qb
		}
	}

	// Regular IN clause with simple values
	qb.where = append(qb.where, &types.WhereCondition{
		Column:    column,
		Operator:  operator,
		Value:     values,
		Connector: connector,
	})

	for _, val := range values {
		qb.bindings = append(qb.bindings, &types.Binding{Value: val})
	}

	return qb
}

// WhereIn adds a WHERE IN condition
func (qb *QueryBuilder) WhereIn(column string, values []any) *QueryBuilder {
	return qb.whereInClause(column, values, "IN", types.AND)
}

// OrWhereIn adds an OR WHERE IN condition
func (qb *QueryBuilder) OrWhereIn(column string, values []any) *QueryBuilder {
	return qb.whereInClause(column, values, "IN", types.OR)
}

// WhereNotIn adds a WHERE NOT IN condition
func (qb *QueryBuilder) WhereNotIn(column string, values []any) *QueryBuilder {
	return qb.whereInClause(column, values, "NOT IN", types.AND)
}

// OrWhereNotIn adds an OR WHERE NOT IN condition
func (qb *QueryBuilder) OrWhereNotIn(column string, values []any) *QueryBuilder {
	return qb.whereInClause(column, values, "NOT IN", types.OR)
}

// WhereInQuery adds a WHERE IN condition with a subquery
func (qb *QueryBuilder) WhereInQuery(column string, sub *QueryBuilder) *QueryBuilder {
	return qb.whereInClause(column, []any{sub}, "IN", types.AND)
}

// OrWhereInQuery adds an OR WHERE IN condition with a subquery
func (qb *QueryBuilder) OrWhereInQuery(column string, sub *QueryBuilder) *QueryBuilder {
	return qb.whereInClause(column, []any{sub}, "IN", types.OR)
}

// WhereNotInQuery adds a WHERE NOT IN condition with a subquery
func (qb *QueryBuilder) WhereNotInQuery(column string, sub *QueryBuilder) *QueryBuilder {
	return qb.whereInClause(column, []any{sub}, "NOT IN", types.AND)
}

// OrWhereNotInQuery adds an OR WHERE NOT IN condition with a subquery
func (qb *QueryBuilder) OrWhereNotInQuery(column string, sub *QueryBuilder) *QueryBuilder {
	return qb.whereInClause(column, []any{sub}, "NOT IN", types.OR)
}

// WhereBoolClause adds a WHERE boolean condition
func (qb *QueryBuilder) whereBoolClause(column string, value bool, operator string, connector types.WhereConditionEnum) *QueryBuilder {
	qb.where = append(qb.where, &types.WhereCondition{
		Column:    column,
		Operator:  operator,
		Value:     value,
		Connector: connector,
	})
	return qb
}

// WhereTrue adds a WHERE true condition
func (qb *QueryBuilder) WhereTrue(column string) *QueryBuilder {
	return qb.whereBoolClause(column, true, "=", types.AND)
}

// OrWhereTrue adds an OR WHERE true condition
func (qb *QueryBuilder) OrWhereTrue(column string) *QueryBuilder {
	return qb.whereBoolClause(column, true, "=", types.OR)
}

// WhereFalse adds a WHERE false condition
func (qb *QueryBuilder) WhereFalse(column string) *QueryBuilder {
	return qb.whereBoolClause(column, false, "=", types.AND)
}

// OrWhereFalse adds an OR WHERE false condition
func (qb *QueryBuilder) OrWhereFalse(column string) *QueryBuilder {
	return qb.whereBoolClause(column, false, "=", types.OR)
}

func (qb *QueryBuilder) whereBetweenClause(column string, min, max any, operator string, connector types.WhereConditionEnum) *QueryBuilder {
	// Support expressions for min/max
	if minExpr, ok := min.(*types.Expression); ok {
		if maxExpr, ok := max.(*types.Expression); ok {
			combined := fmt.Sprintf("%s %s %s %s %s", column, operator, minExpr.Sql, connector, maxExpr.Sql)
			qb.where = append(qb.where, &types.WhereCondition{
				Column:    column,
				Operator:  operator,
				Value:     nil,
				Connector: connector,
				Raw: &types.Expression{
					Sql:      combined,
					Bindings: append(minExpr.Bindings, maxExpr.Bindings...),
				},
			})
			return qb
		}
	}

	qb.where = append(qb.where, &types.WhereCondition{
		Column:    column,
		Operator:  operator,
		Value:     []any{min, max},
		Connector: connector,
	})
	qb.bindings = append(qb.bindings, &types.Binding{Value: min}, &types.Binding{Value: max})
	return qb
}

// WhereBetween adds a WHERE BETWEEN condition
func (qb *QueryBuilder) WhereBetween(column string, min, max any) *QueryBuilder {
	return qb.whereBetweenClause(column, min, max, "BETWEEN", types.AND)
}

// OrWhereBetween adds an OR WHERE BETWEEN condition
func (qb *QueryBuilder) OrWhereBetween(column string, min, max any) *QueryBuilder {
	return qb.whereBetweenClause(column, min, max, "BETWEEN", types.OR)
}

// WhereNotBetween adds a WHERE NOT BETWEEN condition
func (qb *QueryBuilder) WhereNotBetween(column string, min, max any) *QueryBuilder {
	return qb.whereBetweenClause(column, min, max, "NOT BETWEEN", types.AND)
}

// OrWhereNotBetween adds an OR WHERE NOT BETWEEN condition
func (qb *QueryBuilder) OrWhereNotBetween(column string, min, max any) *QueryBuilder {
	return qb.whereBetweenClause(column, min, max, "NOT BETWEEN", types.OR)
}

// WhereExistsClause have the main logic to add a WHERE EXISTS conditions to the query builder
func (qb *QueryBuilder) whereExistsClause(value any, operator string, connector types.WhereConditionEnum) *QueryBuilder {
	if value == nil {
		return qb
	}

	switch v := value.(type) {
	case *types.Expression:
		sqlStr, sqlBindings, err := v.ToSql()
		if err != nil {
			qb.appendError(err)
		}
		qb.where = append(qb.where, &types.WhereCondition{
			Column:    operator,
			Operator:  "",
			Value:     nil,
			Connector: connector,
			Raw: &types.Expression{
				Sql:      operator + " (" + sqlStr + ")",
				Bindings: sqlBindings,
			},
		})
		return qb
	case *types.DialectExpression:
		sqlStr, sqlBindings, err := v.ToSql(qb.GetDialect().GetDriver().String())
		if err != nil {
			qb.appendError(err)
		}
		qb.where = append(qb.where, &types.WhereCondition{
			Column:    operator,
			Operator:  "",
			Value:     nil,
			Connector: connector,
			Raw: &types.Expression{
				Sql:      operator + "(" + sqlStr + ")",
				Bindings: sqlBindings,
			},
		})
		return qb

	case *QueryBuilder:
		subSql, subBindings, err := v.SetDialect(qb.GetDialect().GetDriver()).ToSql()
		if err != nil {
			qb.appendError(err)
		}
		qb.where = append(qb.where, &types.WhereCondition{
			Column:    operator,
			Operator:  "",
			Value:     nil,
			Connector: connector,
			Raw: &types.Expression{
				Sql:      operator + " (" + subSql + ")",
				Bindings: subBindings,
			},
		})
		return qb
	}

	qb.appendError(fmt.Errorf("%w: expected Raw Expression or QueryBuilder in WhereExists clause", xqbErr.ErrInvalidQuery))
	return qb
}

// WhereExists adds a WHERE EXISTS condition it accepts a subquery like Raw Expression or QueryBuilder (implement ToSql)
func (qb *QueryBuilder) WhereExists(subquery any) *QueryBuilder {
	return qb.whereExistsClause(subquery, "EXISTS", types.AND)
}

// WhereExists adds a WHERE EXISTS condition it accepts a subquery like Raw Expression or QueryBuilder (implement ToSql)
func (qb *QueryBuilder) OrWhereExists(subquery any) *QueryBuilder {
	return qb.whereExistsClause(subquery, "EXISTS", types.OR)
}

// WhereNotExists adds a WHERE NOT EXISTS condition it accepts a subquery like Raw Expression or QueryBuilder (implement ToSql)
func (qb *QueryBuilder) WhereNotExists(subquery any) *QueryBuilder {
	return qb.whereExistsClause(subquery, "NOT EXISTS", types.AND)
}

// WhereNotExists adds a WHERE NOT EXISTS condition it accepts a subquery like Raw Expression or QueryBuilder (implement ToSql)
func (qb *QueryBuilder) OrWhereNotExists(subquery any) *QueryBuilder {
	return qb.whereExistsClause(subquery, "NOT EXISTS", types.OR)
}

func (qb *QueryBuilder) whereGroupClause(fn func(qb *QueryBuilder), connector types.WhereConditionEnum) *QueryBuilder {
	// Create a temporary builder to capture the conditions in the group
	groupBuilder := &QueryBuilder{}

	// Execute the function to populate the group builder's conditions
	fn(groupBuilder)

	// If no conditions were added in the group, return the original builder
	if len(groupBuilder.where) == 0 {
		return qb
	}

	// Add the group as a structured Group instead of raw Sql
	qb.where = append(qb.where, &types.WhereCondition{
		Group:     groupBuilder.where,
		Connector: connector,
	})

	return qb
}

// WhereGroup adds a grouped WHERE clause
func (qb *QueryBuilder) WhereGroup(fn func(qb *QueryBuilder)) *QueryBuilder {
	return qb.whereGroupClause(fn, types.AND)
}

// OrWhereGroup adds a grouped WHERE clause with OR connector
func (qb *QueryBuilder) OrWhereGroup(fn func(qb *QueryBuilder)) *QueryBuilder {
	return qb.whereGroupClause(fn, types.OR)
}
