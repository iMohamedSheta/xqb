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
				SQL:      fmt.Sprintf("(%s) %s (%s)", v.SQL, operator, val.SQL),
				Bindings: append(v.Bindings, val.Bindings...),
			}
		default:
			raw = &types.Expression{
				SQL:      fmt.Sprintf("%s %s ?", v.SQL, operator),
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
			subSQL, subBindings, err := v.SetDialect(qb.dialect.GetDriver()).ToSQL()
			if err != nil {
				qb.appendError(err)
			}
			raw = &types.Expression{
				SQL:      fmt.Sprintf("%s %s (%s)", col, operator, subSQL),
				Bindings: subBindings,
			}
		case *types.Expression:
			raw = &types.Expression{
				SQL:      fmt.Sprintf("%s %s (%s)", col, operator, v.SQL),
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
	subSQL, subBindings, _ := sub.SetDialect(qb.dialect.GetDriver()).ToSQL()
	qb.where = append(qb.where, &types.WhereCondition{
		Column:    column,
		Operator:  operator,
		Connector: types.AND,
		Raw: &types.Expression{
			SQL:      fmt.Sprintf("%s %s (%s)", column, operator, subSQL),
			Bindings: subBindings,
		},
	})
	return qb
}

func (qb *QueryBuilder) OrWhereSub(column string, operator string, sub *QueryBuilder) *QueryBuilder {
	subSQL, subBindings, _ := sub.SetDialect(qb.dialect.GetDriver()).ToSQL()
	qb.where = append(qb.where, &types.WhereCondition{
		Column:    column,
		Operator:  operator,
		Connector: types.OR,
		Raw: &types.Expression{
			SQL:      fmt.Sprintf("%s %s (%s)", column, operator, subSQL),
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
			SQL:      fmt.Sprintf("%s %s (%s)", column, operator, expr.SQL),
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
			SQL:      fmt.Sprintf("%s %s (%s)", column, operator, expr.SQL),
			Bindings: expr.Bindings,
		},
	})
	return qb
}

func (qb *QueryBuilder) whereRawClause(sql string, bindings []any, connector types.WhereConditionEnum) *QueryBuilder {
	qb.where = append(qb.where, &types.WhereCondition{
		Raw: &types.Expression{
			SQL:      sql,
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
			subSQL, subBindings, err := v.SetDialect(qb.dialect.GetDriver()).ToSQL()
			if err != nil {
				qb.appendError(err)
			}
			qb.where = append(qb.where, &types.WhereCondition{
				Column:    column,
				Operator:  operator,
				Value:     nil,
				Connector: connector,
				Raw: &types.Expression{
					SQL:      fmt.Sprintf("%s %s (%s)", column, operator, subSQL),
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
					SQL:      fmt.Sprintf("%s %s (%s)", column, operator, v.SQL),
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

func (qb *QueryBuilder) whereBetweenClause(column string, min, max any, operator string, connector types.WhereConditionEnum) *QueryBuilder {
	// Support expressions for min/max
	if minExpr, ok := min.(*types.Expression); ok {
		if maxExpr, ok := max.(*types.Expression); ok {
			combined := fmt.Sprintf("%s %s %s %s %s", column, operator, minExpr.SQL, connector, maxExpr.SQL)
			qb.where = append(qb.where, &types.WhereCondition{
				Column:    column,
				Operator:  operator,
				Value:     nil,
				Connector: connector,
				Raw: &types.Expression{
					SQL:      combined,
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
		sqlStr, sqlBindings, err := v.ToSQL()
		if err != nil {
			qb.appendError(err)
		}
		qb.where = append(qb.where, &types.WhereCondition{
			Column:    operator,
			Operator:  "",
			Value:     nil,
			Connector: connector,
			Raw: &types.Expression{
				SQL:      operator + " (" + sqlStr + ")",
				Bindings: sqlBindings,
			},
		})
		return qb
	case *types.DialectExpression:
		sqlStr, sqlBindings, err := v.ToSQL(qb.dialect.GetDriver().String())
		if err != nil {
			qb.appendError(err)
		}
		qb.where = append(qb.where, &types.WhereCondition{
			Column:    operator,
			Operator:  "",
			Value:     nil,
			Connector: connector,
			Raw: &types.Expression{
				SQL:      operator + "(" + sqlStr + ")",
				Bindings: sqlBindings,
			},
		})
		return qb

	case *QueryBuilder:
		subSQL, subBindings, err := v.SetDialect(qb.dialect.GetDriver()).ToSQL()
		if err != nil {
			qb.appendError(err)
		}
		qb.where = append(qb.where, &types.WhereCondition{
			Column:    operator,
			Operator:  "",
			Value:     nil,
			Connector: connector,
			Raw: &types.Expression{
				SQL:      operator + " (" + subSQL + ")",
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

// func (qb *QueryBuilder) whereGroupClause(fn func(qb *QueryBuilder), connector types.WhereConditionEnum) *QueryBuilder {
// 	// Create a temporary builder to capture the conditions in the group
// 	groupBuilder := &QueryBuilder{}

// 	// Execute the function to populate the group builder's conditions
// 	fn(groupBuilder)

// 	// If no conditions were added in the group, return the original builder
// 	if len(groupBuilder.where) == 0 {
// 		return qb
// 	}

// 	// Create a raw expression for the group
// 	var sql strings.Builder
// 	sql.WriteString("(")

// 	var groupBindings []any

// 	s, bindings, err := groupBuilder.SetDialect(qb.dialect.GetDriver()).ToSQL()
// 	DD(s, bindings, err)
// 	// Process each condition in the group
// 	for i, condition := range groupBuilder.where {
// 		if i > 0 {
// 			sql.WriteString(" ")
// 			sql.WriteString(string(condition.Connector))
// 			sql.WriteString(" ")
// 		}

// 		if condition.Raw != nil {
// 			// Handle raw SQL expression
// 			sql.WriteString(condition.Raw.SQL)
// 			if len(condition.Raw.Bindings) > 0 {
// 				groupBindings = append(groupBindings, condition.Raw.Bindings...)
// 			}
// 		} else {
// 			// Handle regular condition
// 			sql.WriteString(condition.Column)
// 			if condition.Operator != "" {
// 				sql.WriteString(" ")
// 				sql.WriteString(condition.Operator)
// 				if condition.Value != nil {
// 					sql.WriteString(" ?")
// 					groupBindings = append(groupBindings, condition.Value)
// 				}
// 			}
// 		}
// 	}

// 	sql.WriteString(")")

// 	// Add the group as a raw expression to the main builder
// 	qb.where = append(qb.where, &types.WhereCondition{
// 		Raw: &types.Expression{
// 			SQL:      sql.String(),
// 			Bindings: groupBindings,
// 		},
// 		Connector: connector,
// 	})

// 	return qb
// }

func (qb *QueryBuilder) whereGroupClause(fn func(qb *QueryBuilder), connector types.WhereConditionEnum) *QueryBuilder {
	// Create a temporary builder to capture the conditions in the group
	groupBuilder := &QueryBuilder{}

	// Execute the function to populate the group builder's conditions
	fn(groupBuilder)

	// If no conditions were added in the group, return the original builder
	if len(groupBuilder.where) == 0 {
		return qb
	}

	// Add the group as a structured Group instead of raw SQL
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
