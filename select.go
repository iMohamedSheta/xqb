package main

import "github.com/iMohamedSheta/xqb/types"

// Select specifies columns to select
func (qb *QueryBuilder) Select(columns ...any) *QueryBuilder {
	qb.queryType = types.SELECT
	qb.columns = columns
	return qb
}

// SelectRaw adds a raw SQL expression to the SELECT clause
func (qb *QueryBuilder) SelectRaw(sql string, bindings ...interface{}) *QueryBuilder {
	return qb.Select(Raw(sql, bindings...))
}

// AddSelect adds additional columns to select
func (qb *QueryBuilder) AddSelect(columns ...any) *QueryBuilder {
	qb.columns = append(qb.columns, columns...)
	return qb
}

// From specifies the table to select from
func (qb *QueryBuilder) From(table string) *QueryBuilder {
	qb.table = table
	return qb
}

// FromRaw adds a raw SQL expression to the FROM clause
func (qb *QueryBuilder) FromRaw(expression string) *QueryBuilder {
	qb.table = expression
	return qb
}

// FromSubquery uses a subquery as the FROM clause
func (qb *QueryBuilder) FromSubquery(subquery *QueryBuilder, alias string) *QueryBuilder {
	qb.subqueries[alias] = subquery
	return qb
}

// Distinct adds DISTINCT to the query
func (qb *QueryBuilder) Distinct() *QueryBuilder {
	qb.isUsingDistinct = true
	return qb
}

// Alias sets an alias for a column
func (qb *QueryBuilder) Alias(column string, alias string) *QueryBuilder {
	qb.columnAliases[column] = alias
	return qb
}
