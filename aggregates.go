package xqb

import "github.com/iMohamedSheta/xqb/types"

// Aggregate adds an aggregate function
func (qb *QueryBuilder) Aggregate(function types.AggregateFunction, column string, alias string) *QueryBuilder {
	qb.aggregateFuncs = append(qb.aggregateFuncs, types.AggregateExpr{
		Function: function,
		Column:   column,
		Alias:    alias,
	})
	return qb
}

// Count adds a COUNT function
func (qb *QueryBuilder) Count(column string, alias ...string) (string, []interface{}, error) {
	al := column + "_count"
	if len(alias) > 0 {
		al = alias[0]
	}
	qb.aggregateFuncs = append(qb.aggregateFuncs, types.AggregateExpr{
		Function: types.COUNT,
		Column:   column,
		Alias:    al,
	})
	return qb.ToSQL()
}

// Sum adds a SUM function
func (qb *QueryBuilder) Sum(column string, alias ...string) *QueryBuilder {
	al := column + "_sum"
	if len(alias) > 0 {
		al = alias[0]
	}
	qb.aggregateFuncs = append(qb.aggregateFuncs, types.AggregateExpr{
		Function: types.SUM,
		Column:   column,
		Alias:    al,
	})
	return qb
}

// Avg adds an AVG function
func (qb *QueryBuilder) Avg(column string, alias ...string) *QueryBuilder {
	al := column + "_avg"
	if len(alias) > 0 {
		al = alias[0]
	}
	qb.aggregateFuncs = append(qb.aggregateFuncs, types.AggregateExpr{
		Function: types.AVG,
		Column:   column,
		Alias:    al,
	})
	return qb
}

// Min adds a MIN function
func (qb *QueryBuilder) Min(column string, alias ...string) *QueryBuilder {
	al := column + "_min"
	if len(alias) > 0 {
		al = alias[0]
	}
	qb.aggregateFuncs = append(qb.aggregateFuncs, types.AggregateExpr{
		Function: types.MIN,
		Column:   column,
		Alias:    al,
	})
	return qb
}

// Max adds a MAX function
func (qb *QueryBuilder) Max(column string, alias ...string) *QueryBuilder {
	al := column + "_max"
	if len(alias) > 0 {
		al = alias[0]
	}
	qb.aggregateFuncs = append(qb.aggregateFuncs, types.AggregateExpr{
		Function: types.MAX,
		Column:   column,
		Alias:    al,
	})
	return qb
}

// JSON adds a JSON expression
func (qb *QueryBuilder) JSON(column string, path string, alias string) *QueryBuilder {
	qb.jsonExpressions = append(qb.jsonExpressions, types.JSONExpression{
		Column: column,
		Path:   path,
		Alias:  alias,
	})
	return qb
}

// JSONExtract adds a JSON_EXTRACT function
func (qb *QueryBuilder) JSONExtract(column string, path string, alias string) *QueryBuilder {
	qb.jsonExpressions = append(qb.jsonExpressions, types.JSONExpression{
		Column:   column,
		Path:     path,
		Function: "JSON_EXTRACT",
		Alias:    alias,
	})
	return qb
}

// Math adds a mathematical expression
func (qb *QueryBuilder) Math(expression string, alias string) *QueryBuilder {
	qb.mathExpressions = append(qb.mathExpressions, types.MathExpression{
		Expression: expression,
		Alias:      alias,
	})
	return qb
}

// Conditional adds a conditional expression
func (qb *QueryBuilder) Conditional(expression string, alias string) *QueryBuilder {
	qb.conditionalExprs = append(qb.conditionalExprs, types.ConditionalExpr{
		Expression: expression,
		Alias:      alias,
	})
	return qb
}

// String adds a string function
func (qb *QueryBuilder) String(function string, column string, params []interface{}, alias string) *QueryBuilder {
	qb.stringFuncs = append(qb.stringFuncs, types.StringFunction{
		Function: function,
		Column:   column,
		Params:   params,
		Alias:    alias,
	})
	return qb
}

// Concat adds a CONCAT function
func (qb *QueryBuilder) Concat(columns []string, separator string, alias string) *QueryBuilder {
	params := make([]interface{}, len(columns)+1)
	params[0] = separator
	for i, col := range columns {
		params[i+1] = col
	}
	qb.stringFuncs = append(qb.stringFuncs, types.StringFunction{
		Function: "CONCAT_WS",
		Column:   "",
		Params:   params,
		Alias:    alias,
	})
	return qb
}

// Date adds a date function
func (qb *QueryBuilder) Date(function string, column string, params []interface{}, alias string) *QueryBuilder {
	qb.dateFuncs = append(qb.dateFuncs, types.DateFunction{
		Function: function,
		Column:   column,
		Params:   params,
		Alias:    alias,
	})
	return qb
}

// DateFormat adds a DATE_FORMAT function
func (qb *QueryBuilder) DateFormat(column string, format string, alias string) *QueryBuilder {
	qb.dateFuncs = append(qb.dateFuncs, types.DateFunction{
		Function: "DATE_FORMAT",
		Column:   column,
		Params:   []interface{}{format},
		Alias:    alias,
	})
	return qb
}

// IndexHint adds an index hint
func (qb *QueryBuilder) IndexHint(hint string) *QueryBuilder {
	qb.indexHints = append(qb.indexHints, hint)
	return qb
}

// ForceIndex adds a FORCE INDEX hint
func (qb *QueryBuilder) ForceIndex(indexName string) *QueryBuilder {
	qb.forceIndex = indexName
	return qb
}

// UseIndex adds a USE INDEX hint
func (qb *QueryBuilder) UseIndex(indexName string) *QueryBuilder {
	qb.useIndex = indexName
	return qb
}

// IgnoreIndex adds an IGNORE INDEX hint
func (qb *QueryBuilder) IgnoreIndex(indexName string) *QueryBuilder {
	qb.ignoreIndex = indexName
	return qb
}

// Lock sets a FOR UPDATE lock
func (qb *QueryBuilder) ForUpdate() *QueryBuilder {
	qb.isForUpdate = true
	return qb
}

// LockInShareMode sets a LOCK IN SHARE MODE lock
func (qb *QueryBuilder) LockInShareMode() *QueryBuilder {
	qb.isLockInShareMode = true
	return qb
}

// HighPriority adds HIGH_PRIORITY hint
func (qb *QueryBuilder) HighPriority() *QueryBuilder {
	qb.isHighPriority = true
	return qb
}

// StraightJoin adds STRAIGHT_JOIN hint
func (qb *QueryBuilder) StraightJoin() *QueryBuilder {
	qb.isStraightJoin = true
	return qb
}

// CalcFoundRows adds SQL_CALC_FOUND_ROWS hint
func (qb *QueryBuilder) CalcFoundRows() *QueryBuilder {
	qb.isCalcFoundRows = true
	return qb
}

// Procedure adds a PROCEDURE call
func (qb *QueryBuilder) Procedure(procedure string, params ...interface{}) *QueryBuilder {
	qb.procedure = procedure
	qb.procedureParams = params
	return qb
}

// Comment adds a comment to the query
func (qb *QueryBuilder) Comment(comment string) *QueryBuilder {
	qb.comment = comment
	return qb
}
