package xqb

import (
	"fmt"

	"github.com/iMohamedSheta/xqb/shared/types"
)

func (qb *QueryBuilder) addJoin(joinType types.JoinType, table any, condition any, alias string, values ...any) *QueryBuilder {
	var tableSQL string
	var conditionSQL string
	var bindings []types.Binding

	// Handle table
	switch t := table.(type) {
	case string:
		tableSQL = t
	case *QueryBuilder:
		subSQL, subBindings, _ := t.ToSQL()
		if alias == "" {
			alias = "sub" // fallback alias
		}
		tableSQL = fmt.Sprintf("(%s) AS %s", subSQL, alias)
		for _, b := range subBindings {
			bindings = append(bindings, types.Binding{Value: b})
		}
	case *types.Expression:
		tableSQL = t.SQL
		for _, b := range t.Bindings {
			bindings = append(bindings, types.Binding{Value: b})
		}
	}

	// Handle condition
	switch c := condition.(type) {
	case string:
		conditionSQL = c
		for _, val := range values {
			bindings = append(bindings, types.Binding{Value: val})
		}
	case *types.Expression:
		conditionSQL = c.SQL
		for _, b := range c.Bindings {
			bindings = append(bindings, types.Binding{Value: b})
		}
	}

	qb.joins = append(qb.joins, &types.Join{
		Type:      joinType,
		Table:     tableSQL,
		Condition: conditionSQL,
		Binding:   bindings,
	})

	return qb
}

// Join adds a INNER JOIN clause to the query
func (qb *QueryBuilder) Join(table string, condition string, values ...any) *QueryBuilder {
	return qb.addJoin(types.INNER_JOIN, table, condition, "", values...)
}

// LeftJoin adds a LEFT JOIN clause to the query
func (qb *QueryBuilder) LeftJoin(table string, condition string, values ...any) *QueryBuilder {
	return qb.addJoin(types.LEFT_JOIN, table, condition, "", values...)
}

// RightJoin adds a RIGHT JOIN clause to the query
func (qb *QueryBuilder) RightJoin(table string, condition string, values ...any) *QueryBuilder {
	return qb.addJoin(types.RIGHT_JOIN, table, condition, "", values...)
}

// FullJoin adds a FULL JOIN clause to the query
func (qb *QueryBuilder) FullJoin(table string, condition string, values ...any) *QueryBuilder {
	return qb.addJoin(types.FULL_JOIN, table, condition, "", values...)
}

// CrossJoin adds a CROSS JOIN clause to the query
func (qb *QueryBuilder) CrossJoin(table string) *QueryBuilder {
	return qb.addJoin(types.CROSS_JOIN, table, "", "")
}

// JoinSub adds a JOIN clause with a subquery
func (qb *QueryBuilder) JoinSub(sub *QueryBuilder, alias, condition string, values ...any) *QueryBuilder {
	return qb.addJoin(types.INNER_JOIN, sub, condition, alias, values...)
}

// LeftJoinSub adds a LEFT JOIN clause with a subquery
func (qb *QueryBuilder) LeftJoinSub(sub *QueryBuilder, alias, condition string, values ...any) *QueryBuilder {
	return qb.addJoin(types.LEFT_JOIN, sub, condition, alias, values...)
}

// RightJoinSub adds a RIGHT JOIN clause with a subquery
func (qb *QueryBuilder) RightJoinSub(sub *QueryBuilder, alias, condition string, values ...any) *QueryBuilder {
	return qb.addJoin(types.RIGHT_JOIN, sub, condition, alias, values...)
}

// FullJoinSub adds a FULL JOIN clause with a subquery
func (qb *QueryBuilder) FullJoinSub(sub *QueryBuilder, alias, condition string, values ...any) *QueryBuilder {
	return qb.addJoin(types.FULL_JOIN, sub, condition, alias, values...)
}

// CrossJoinSub adds a CROSS JOIN clause with a subquery
func (qb *QueryBuilder) CrossJoinSub(sub *QueryBuilder, alias string) *QueryBuilder {
	return qb.addJoin(types.CROSS_JOIN, sub, "", alias)
}

// JoinExpr adds a JOIN clause with an expression
func (qb *QueryBuilder) JoinExpr(expr *types.Expression, condition any, values ...any) *QueryBuilder {
	return qb.addJoin(types.INNER_JOIN, expr, condition, "", values...)
}

// LeftJoinExpr adds a LEFT JOIN clause with an expression
func (qb *QueryBuilder) LeftJoinExpr(expr *types.Expression, condition any, values ...any) *QueryBuilder {
	return qb.addJoin(types.LEFT_JOIN, expr, condition, "", values...)
}

// RightJoinExpr adds a RIGHT JOIN clause with an expression
func (qb *QueryBuilder) RightJoinExpr(expr *types.Expression, condition any, values ...any) *QueryBuilder {
	return qb.addJoin(types.RIGHT_JOIN, expr, condition, "", values...)
}

// FullJoinExpr adds a FULL JOIN clause with an expression
func (qb *QueryBuilder) FullJoinExpr(expr *types.Expression, condition any, values ...any) *QueryBuilder {
	return qb.addJoin(types.FULL_JOIN, expr, condition, "", values...)
}

// CrossJoinExpr adds a CROSS JOIN clause with an expression
func (qb *QueryBuilder) CrossJoinExpr(expr *types.Expression) *QueryBuilder {
	return qb.addJoin(types.CROSS_JOIN, expr, "", "")
}
