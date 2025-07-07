package xqb

import (
	"fmt"
	"strings"

	"github.com/iMohamedSheta/xqb/types"
)

// Where adds a WHERE condition
func (qb *QueryBuilder) Where(column any, operator string, value interface{}) *QueryBuilder {
	var col string
	var bindings []interface{}

	switch v := column.(type) {
	case string:
		col = v
	case *types.Expression:
		col = v.SQL
		bindings = v.Bindings
		// If a value is provided with a raw expression, add it to bindings
		if value != nil {
			bindings = append(bindings, value)
		}
	default:
		col = fmt.Sprintf("%v", v)
	}

	qb.where = append(qb.where, types.WhereCondition{
		Column:    col,
		Operator:  operator,
		Value:     value,
		Connector: types.AND,
		Raw:       nil,
	})

	// Add all bindings
	if len(bindings) > 0 {
		for _, binding := range bindings {
			qb.bindings = append(qb.bindings, types.Binding{Value: binding})
		}
	} else if value != nil {
		qb.bindings = append(qb.bindings, types.Binding{Value: value})
	}

	return qb
}

// WhereRaw adds a raw WHERE condition
func (qb *QueryBuilder) WhereRaw(sql string, bindings ...interface{}) *QueryBuilder {
	expr := Raw(sql, bindings...)
	qb.where = append(qb.where, types.WhereCondition{
		Column:    expr.SQL,
		Operator:  "",
		Value:     nil,
		Connector: types.AND,
		Raw:       expr,
	})
	return qb
}

// OrWhere adds an OR WHERE condition
func (qb *QueryBuilder) OrWhere(column any, operator string, value interface{}) *QueryBuilder {
	var col string
	var bindings []interface{}

	switch v := column.(type) {
	case string:
		col = v
	case *types.Expression:
		col = v.SQL
		bindings = v.Bindings
	default:
		col = fmt.Sprintf("%v", v)
	}

	qb.where = append(qb.where, types.WhereCondition{
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

// OrWhereRaw adds a raw OR WHERE condition
func (qb *QueryBuilder) OrWhereRaw(sql string, bindings ...interface{}) *QueryBuilder {
	return qb.OrWhere(Raw(sql, bindings...), "", nil)
}

// WhereNull adds a WHERE NULL condition
func (qb *QueryBuilder) WhereNull(column string) *QueryBuilder {
	qb.where = append(qb.where, types.WhereCondition{
		Column:    column,
		Operator:  "IS NULL",
		Value:     nil,
		Connector: types.AND,
	})
	return qb
}

// WhereNotNull adds a WHERE NOT NULL condition
func (qb *QueryBuilder) WhereNotNull(column string) *QueryBuilder {
	qb.where = append(qb.where, types.WhereCondition{
		Column:    column,
		Operator:  "IS NOT NULL",
		Value:     nil,
		Connector: types.AND,
	})
	return qb
}

// WhereIn adds a WHERE IN condition
func (qb *QueryBuilder) WhereIn(column string, values []interface{}) *QueryBuilder {
	qb.where = append(qb.where, types.WhereCondition{
		Column:    column,
		Operator:  "IN",
		Value:     values,
		Connector: types.AND,
	})
	return qb
}

// WhereNotIn adds a WHERE NOT IN condition
func (qb *QueryBuilder) WhereNotIn(column string, values []interface{}) *QueryBuilder {
	qb.where = append(qb.where, types.WhereCondition{
		Column:    column,
		Operator:  "NOT IN",
		Value:     values,
		Connector: types.AND,
		IsNot:     true,
	})
	return qb
}

// WhereBetween adds a WHERE BETWEEN condition
func (qb *QueryBuilder) WhereBetween(column string, min, max interface{}) *QueryBuilder {
	qb.where = append(qb.where, types.WhereCondition{
		Column:    column,
		Operator:  "BETWEEN",
		Value:     []interface{}{min, max},
		Connector: types.AND,
	})
	return qb
}

// WhereNotBetween adds a WHERE NOT BETWEEN condition
func (qb *QueryBuilder) WhereNotBetween(column string, min, max interface{}) *QueryBuilder {
	qb.where = append(qb.where, types.WhereCondition{
		Column:    column,
		Operator:  "NOT BETWEEN",
		Value:     []interface{}{min, max},
		Connector: types.AND,
		IsNot:     true,
	})
	return qb
}

// WhereExists adds a WHERE EXISTS condition
func (qb *QueryBuilder) WhereExists(subquery *QueryBuilder) *QueryBuilder {
	// Get the SQL and bindings from the subquery
	subSQL, subBindings, _ := subquery.ToSQL()

	// Create a raw EXISTS expression with the subquery
	qb.where = append(qb.where, types.WhereCondition{
		Column:    "EXISTS",
		Operator:  "",
		Value:     nil,
		Connector: types.AND,
		Raw: &types.Expression{
			SQL:      "EXISTS (" + subSQL + ")",
			Bindings: subBindings,
		},
	})
	return qb
}

// WhereNotExists adds a WHERE NOT EXISTS condition
func (qb *QueryBuilder) WhereNotExists(subquery *QueryBuilder) *QueryBuilder {
	// Generate a unique alias for the subquery
	alias := fmt.Sprintf("subquery_%d", len(qb.subqueries))
	qb.subqueries[alias] = subquery

	qb.where = append(qb.where, types.WhereCondition{
		Column:    "NOT EXISTS",
		Operator:  "",
		Value:     alias,
		Connector: types.AND,
		IsNot:     true,
	})
	return qb
}

// WhereGroup adds a grouped WHERE clause
func (qb *QueryBuilder) WhereGroup(fn func(qb *QueryBuilder)) *QueryBuilder {
	// Create a temporary builder to capture the conditions in the group
	groupBuilder := &QueryBuilder{
		where: []types.WhereCondition{},
	}

	// Execute the function to populate the group builder's conditions
	fn(groupBuilder)

	// If no conditions were added in the group, return the original builder
	if len(groupBuilder.where) == 0 {
		return qb
	}

	// Create a raw expression for the group
	var sql strings.Builder
	sql.WriteString("(")

	var groupBindings []interface{}

	// Process each condition in the group
	for i, condition := range groupBuilder.where {
		if i > 0 {
			sql.WriteString(" ")
			sql.WriteString(string(condition.Connector))
			sql.WriteString(" ")
		}

		if condition.Raw != nil {
			// Handle raw SQL expression
			sql.WriteString(condition.Raw.SQL)
			if len(condition.Raw.Bindings) > 0 {
				groupBindings = append(groupBindings, condition.Raw.Bindings...)
			}
		} else {
			// Handle regular condition
			sql.WriteString(condition.Column)
			if condition.Operator != "" {
				sql.WriteString(" ")
				sql.WriteString(condition.Operator)
				if condition.Value != nil {
					sql.WriteString(" ?")
					groupBindings = append(groupBindings, condition.Value)
				}
			}
		}
	}

	sql.WriteString(")")

	// Add the group as a raw expression to the main builder
	qb.where = append(qb.where, types.WhereCondition{
		Raw: &types.Expression{
			SQL:      sql.String(),
			Bindings: groupBindings,
		},
		Connector: types.AND,
	})

	return qb
}

// OrWhereGroup adds a grouped WHERE clause with OR connector
func (qb *QueryBuilder) OrWhereGroup(fn func(qb *QueryBuilder)) *QueryBuilder {
	// Create a temporary builder to capture the conditions in the group
	groupBuilder := &QueryBuilder{
		where: []types.WhereCondition{},
	}

	// Execute the function to populate the group builder's conditions
	fn(groupBuilder)

	// If no conditions were added in the group, return the original builder
	if len(groupBuilder.where) == 0 {
		return qb
	}

	// Create a raw expression for the group
	var sql strings.Builder
	sql.WriteString("(")

	var groupBindings []interface{}

	// Process each condition in the group
	for i, condition := range groupBuilder.where {
		if i > 0 {
			sql.WriteString(" ")
			sql.WriteString(string(condition.Connector))
			sql.WriteString(" ")
		}

		if condition.Raw != nil {
			// Handle raw SQL expression
			sql.WriteString(condition.Raw.SQL)
			if len(condition.Raw.Bindings) > 0 {
				groupBindings = append(groupBindings, condition.Raw.Bindings...)
			}
		} else {
			// Handle regular condition
			sql.WriteString(condition.Column)
			if condition.Operator != "" {
				sql.WriteString(" ")
				sql.WriteString(condition.Operator)
				if condition.Value != nil {
					sql.WriteString(" ?")
					groupBindings = append(groupBindings, condition.Value)
				}
			}
		}
	}

	sql.WriteString(")")

	// Add the group as a raw expression to the main builder with OR connector
	qb.where = append(qb.where, types.WhereCondition{
		Raw: &types.Expression{
			SQL:      sql.String(),
			Bindings: groupBindings,
		},
		Connector: types.OR,
	})

	return qb
}
