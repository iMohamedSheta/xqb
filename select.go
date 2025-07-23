package xqb

import (
	"github.com/iMohamedSheta/xqb/shared/enums"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// Select specifies columns to select
func (qb *QueryBuilder) Select(columns ...any) *QueryBuilder {
	qb.queryType = enums.SELECT
	qb.columns = columns
	return qb
}

// SelectSub add sub queries to select
func (qb *QueryBuilder) SelectSub(subQuery *QueryBuilder, alias string) *QueryBuilder {
	sql, bindings, err := subQuery.SetDialect(qb.dialect.GetDriver()).ToSQL()
	if err != nil {
		qb.errors = append(qb.errors, err)
		return qb
	}
	qb.queryType = enums.SELECT
	qb.columns = append(qb.columns, Raw("("+sql+") AS "+alias, bindings...))
	return qb
}

// SelectRaw adds a raw SQL expression to the SELECT clause
func (qb *QueryBuilder) SelectRaw(sql string, bindings ...any) *QueryBuilder {
	return qb.Select(Raw(sql, bindings...))
}

// AddSelect adds additional columns to select
func (qb *QueryBuilder) AddSelect(columns ...any) *QueryBuilder {
	qb.columns = append(qb.columns, columns...)
	return qb
}

// SelectRaw adds a raw SQL expression to the SELECT clause
func (qb *QueryBuilder) AddSelectRaw(sql string, bindings ...any) *QueryBuilder {
	qb.AddSelect(Raw(sql, bindings...))
	return qb
}

// From specifies the table to select from
func (qb *QueryBuilder) From(table string) *QueryBuilder {
	qb.table = &types.Table{
		Name: table,
	}
	return qb
}

// FromSubquery uses a subquery as the FROM clause
func (qb *QueryBuilder) FromSubquery(subQuery *QueryBuilder, alias string) *QueryBuilder {
	raw := subQuery.SetDialect(qb.dialect.GetDriver()).ToRawExpr()
	raw.SQL = "(" + raw.SQL + ")" + " AS " + alias
	qb.table = &types.Table{Raw: raw}
	return qb
}

// Distinct adds DISTINCT to the query
func (qb *QueryBuilder) Distinct() *QueryBuilder {
	qb.isUsingDistinct = true
	return qb
}
