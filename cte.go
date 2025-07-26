package xqb

import (
	"github.com/iMohamedSheta/xqb/shared/types"
)

func (qb *QueryBuilder) With(name string, query *QueryBuilder) *QueryBuilder {
	return qb.with(name, query.GetData(), false)
}

func (qb *QueryBuilder) WithRaw(name string, sql string, values ...any) *QueryBuilder {
	return qb.WithExpr(name, sql, values...)
}

func (qb *QueryBuilder) WithExpr(name string, sql string, bindings ...any) *QueryBuilder {
	return qb.withExpr(name, sql, false, bindings...)
}

func (qb *QueryBuilder) WithRecursive(name string, query *QueryBuilder) *QueryBuilder {
	return qb.with(name, query.GetData(), true)
}

func (qb *QueryBuilder) WithRecursiveRaw(name string, sql string, values ...any) *QueryBuilder {
	return qb.withExpr(name, sql, true, values...)
}

func (qb *QueryBuilder) with(name string, data *types.QueryBuilderData, recursive bool) *QueryBuilder {
	qb.withCTEs = append(qb.withCTEs, &types.CTE{
		Name:      name,
		Query:     data,
		Recursive: recursive,
	})
	return qb
}

func (qb *QueryBuilder) withExpr(name string, sql string, recursive bool, bindings ...any) *QueryBuilder {
	qb.withCTEs = append(qb.withCTEs, &types.CTE{
		Name: name,
		Expression: &types.Expression{
			Sql:      sql,
			Bindings: bindings,
		},
		Recursive: recursive,
	})
	return qb
}
