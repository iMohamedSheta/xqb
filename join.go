package main

import (
	"github.com/iMohamedSheta/xqb/types"
)

// Join adds a JOIN clause Example: qb.Join("users", "users.id = posts.user_id AND users.name = ?", "John")
func (qb *QueryBuilder) Join(table string, condition string, values ...interface{}) *QueryBuilder {
	join := types.Join{
		Type:      types.INNER_JOIN,
		Table:     table,
		Condition: condition,
		Binding:   []types.Binding{},
	}

	// Add all bindings
	if len(values) > 0 {
		for _, value := range values {
			join.Binding = append(join.Binding, types.Binding{Value: value})
		}
	}

	qb.joins = append(qb.joins, join)
	return qb
}

// LeftJoin adds a LEFT JOIN clause
func (qb *QueryBuilder) LeftJoin(table string, condition string, values ...interface{}) *QueryBuilder {
	join := types.Join{
		Type:      types.LEFT_JOIN,
		Table:     table,
		Condition: condition,
		Binding:   []types.Binding{},
	}

	// Add all bindings
	if len(values) > 0 {
		for _, value := range values {
			join.Binding = append(join.Binding, types.Binding{Value: value})
		}
	}

	qb.joins = append(qb.joins, join)

	return qb
}

// RightJoin adds a RIGHT JOIN clause
func (qb *QueryBuilder) RightJoin(table string, condition string, values ...interface{}) *QueryBuilder {
	join := types.Join{
		Type:      types.RIGHT_JOIN,
		Table:     table,
		Condition: condition,
		Binding:   []types.Binding{},
	}

	// Add all bindings
	if len(values) > 0 {
		for _, value := range values {
			join.Binding = append(join.Binding, types.Binding{Value: value})
		}
	}

	qb.joins = append(qb.joins, join)

	return qb
}

// CrossJoin adds a CROSS JOIN clause
func (qb *QueryBuilder) CrossJoin(table string) *QueryBuilder {
	qb.joins = append(qb.joins, types.Join{
		Type:  types.CROSS_JOIN,
		Table: table,
	})
	return qb
}
