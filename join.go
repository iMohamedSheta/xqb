package xqb

import (
	"fmt"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

func (qb *QueryBuilder) addJoin(joinType types.JoinType, table any, condition any, alias string, values ...any) *QueryBuilder {
	var tableSql string
	var bindings []types.Binding

	dialect := qb.GetDialect().Getdialect().String()

	// Handle table
	switch t := table.(type) {
	case string:
		tableSql = t
	case *QueryBuilder:
		subSql, subBindings, err := t.SetDialect(qb.GetDialect().Getdialect()).ToSql()
		if err != nil {
			qb.appendError(err)
		}
		if alias == "" {
			qb.appendError(fmt.Errorf("%w: alias is required for subquery or expression join", xqbErr.ErrInvalidQuery))
		}
		tableSql = fmt.Sprintf("(%s) AS %s", subSql, alias)
		for _, b := range subBindings {
			bindings = append(bindings, types.Binding{Value: b})
		}
	case types.ExpressionInterface:
		exprSql, exprBindings, err := t.ToSql(dialect)
		if err != nil {
			qb.appendError(err)
			return qb
		}
		tableSql = exprSql
		for _, b := range exprBindings {
			bindings = append(bindings, types.Binding{Value: b})
		}
	}

	// Handle condition — three variants:
	// string:           "users.id = orders.user_id" with optional ? values
	// func(*JoinClause): structured On/Where builder (preferred)
	// ExpressionInterface: Raw/DialectExpression
	switch c := condition.(type) {
	case func(*JoinClause):
		// closure-based: compiler extracts SQL + bindings from Conditions
		clause := &JoinClause{}
		c(clause)
		qb.joins = append(qb.joins, &types.Join{
			Type:       joinType,
			Table:      tableSql,
			Conditions: clause.conditions,
			Binding:    bindings,
		})
		return qb

	case string:
		// plain string: Join("orders", "users.id = orders.user_id AND orders.type = ?", "invoice")
		for _, val := range values {
			bindings = append(bindings, types.Binding{Value: val})
		}
		qb.joins = append(qb.joins, &types.Join{
			Type:      joinType,
			Table:     tableSql,
			Condition: c,
			Binding:   bindings,
		})
		return qb

	case types.ExpressionInterface:
		exprSql, exprBindings, err := c.ToSql(dialect)
		if err != nil {
			qb.appendError(err)
			return qb
		}
		for _, b := range exprBindings {
			bindings = append(bindings, types.Binding{Value: b})
		}
		qb.joins = append(qb.joins, &types.Join{
			Type:      joinType,
			Table:     tableSql,
			Condition: exprSql,
			Binding:   bindings,
		})
		return qb
	}

	qb.appendError(fmt.Errorf("%w: unsupported condition type for JOIN", xqbErr.ErrInvalidQuery))
	return qb
}

// Join adds a INNER JOIN clause to the query
func (qb *QueryBuilder) Join(table string, condition any, values ...any) *QueryBuilder {
	return qb.addJoin(types.INNER_JOIN, table, condition, "", values...)
}

// LeftJoin adds a LEFT JOIN clause to the query
func (qb *QueryBuilder) LeftJoin(table string, condition any, values ...any) *QueryBuilder {
	return qb.addJoin(types.LEFT_JOIN, table, condition, "", values...)
}

// RightJoin adds a RIGHT JOIN clause to the query
func (qb *QueryBuilder) RightJoin(table string, condition any, values ...any) *QueryBuilder {
	return qb.addJoin(types.RIGHT_JOIN, table, condition, "", values...)
}

// FullJoin adds a FULL JOIN clause to the query
func (qb *QueryBuilder) FullJoin(table string, condition any, values ...any) *QueryBuilder {
	return qb.addJoin(types.FULL_JOIN, table, condition, "", values...)
}

// CrossJoin adds a CROSS JOIN clause to the query
func (qb *QueryBuilder) CrossJoin(table string) *QueryBuilder {
	return qb.addJoin(types.CROSS_JOIN, table, "", "")
}

// JoinSub adds a JOIN clause with a subquery
func (qb *QueryBuilder) JoinSub(sub *QueryBuilder, alias string, condition any, values ...any) *QueryBuilder {
	return qb.addJoin(types.INNER_JOIN, sub, condition, alias, values...)
}

// LeftJoinSub adds a LEFT JOIN clause with a subquery
func (qb *QueryBuilder) LeftJoinSub(sub *QueryBuilder, alias string, condition any, values ...any) *QueryBuilder {
	return qb.addJoin(types.LEFT_JOIN, sub, condition, alias, values...)
}

// RightJoinSub adds a RIGHT JOIN clause with a subquery
func (qb *QueryBuilder) RightJoinSub(sub *QueryBuilder, alias string, condition any, values ...any) *QueryBuilder {
	return qb.addJoin(types.RIGHT_JOIN, sub, condition, alias, values...)
}

// FullJoinSub adds a FULL JOIN clause with a subquery
func (qb *QueryBuilder) FullJoinSub(sub *QueryBuilder, alias string, condition any, values ...any) *QueryBuilder {
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
