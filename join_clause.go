package xqb

import "github.com/iMohamedSheta/xqb/shared/types"

// JoinClause builds structured ON conditions for a JOIN.
// Obtain one via the closure passed to Join/LeftJoin/etc.
type JoinClause struct {
	conditions []*types.JoinCondition
}

// On adds an AND ON condition comparing two columns (no binding)
// e.g. On("users.id", "=", "orders.user_id")
func (j *JoinClause) On(first, operator, second string) *JoinClause {
	return j.onClause(first, operator, second, types.AND)
}

// OrOn adds an OR ON condition comparing two columns (no binding)
func (j *JoinClause) OrOn(first, operator, second string) *JoinClause {
	return j.onClause(first, operator, second, types.OR)
}

func (j *JoinClause) onClause(first, operator, second string, connector types.WhereConditionEnum) *JoinClause {
	j.conditions = append(j.conditions, &types.JoinCondition{
		Kind:      types.JoinConditionOn,
		First:     first,
		Operator:  operator,
		Second:    second,
		Connector: connector,
	})
	return j
}

// Where adds an AND ON condition comparing a column to a bound value
// e.g. Where("orders.type", "=", "invoice")
func (j *JoinClause) Where(column, operator string, value any) *JoinClause {
	return j.whereClause(column, operator, value, types.AND)
}

// OrWhere adds an OR ON condition comparing a column to a bound value
func (j *JoinClause) OrWhere(column, operator string, value any) *JoinClause {
	return j.whereClause(column, operator, value, types.OR)
}

func (j *JoinClause) whereClause(column, operator string, value any, connector types.WhereConditionEnum) *JoinClause {
	j.conditions = append(j.conditions, &types.JoinCondition{
		Kind:      types.JoinConditionWhere,
		First:     column,
		Operator:  operator,
		Value:     value,
		Connector: connector,
	})
	return j
}

// WhereNull adds an AND ON IS NULL condition
func (j *JoinClause) WhereNull(column string) *JoinClause {
	return j.whereClause(column, "IS NULL", nil, types.AND)
}

// OrWhereNull adds an OR ON IS NULL condition
func (j *JoinClause) OrWhereNull(column string) *JoinClause {
	return j.whereClause(column, "IS NULL", nil, types.OR)
}

// WhereNotNull adds an AND ON IS NOT NULL condition
func (j *JoinClause) WhereNotNull(column string) *JoinClause {
	return j.whereClause(column, "IS NOT NULL", nil, types.AND)
}

// OrWhereNotNull adds an OR ON IS NOT NULL condition
func (j *JoinClause) OrWhereNotNull(column string) *JoinClause {
	return j.whereClause(column, "IS NOT NULL", nil, types.OR)
}

// OnRaw adds a raw AND ON SQL fragment
// e.g. OnRaw("users.id = orders.user_id AND orders.active = ?", 1)
func (j *JoinClause) OnRaw(sql string, bindings ...any) *JoinClause {
	return j.onRawClause(sql, bindings, types.AND)
}

// OrOnRaw adds a raw OR ON SQL fragment
func (j *JoinClause) OrOnRaw(sql string, bindings ...any) *JoinClause {
	return j.onRawClause(sql, bindings, types.OR)
}

func (j *JoinClause) onRawClause(sql string, bindings []any, connector types.WhereConditionEnum) *JoinClause {
	j.conditions = append(j.conditions, &types.JoinCondition{
		Kind:      types.JoinConditionRaw,
		Connector: connector,
		Raw:       &types.Expression{Sql: sql, Bindings: bindings},
	})
	return j
}

// OnGroup adds a nested AND group of ON conditions
// e.g. OnGroup(func(j *JoinClause) { j.On(...).OrOn(...) })
func (j *JoinClause) OnGroup(fn func(*JoinClause)) *JoinClause {
	return j.onGroupClause(fn, types.AND)
}

// OrOnGroup adds a nested OR group of ON conditions
func (j *JoinClause) OrOnGroup(fn func(*JoinClause)) *JoinClause {
	return j.onGroupClause(fn, types.OR)
}

func (j *JoinClause) onGroupClause(fn func(*JoinClause), connector types.WhereConditionEnum) *JoinClause {
	group := &JoinClause{}
	fn(group)
	if len(group.conditions) == 0 {
		return j
	}
	j.conditions = append(j.conditions, &types.JoinCondition{
		Kind:      types.JoinConditionGroup,
		Connector: connector,
		Group:     group.conditions,
	})
	return j
}
