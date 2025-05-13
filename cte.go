package main

import (
	"github.com/iMohamedSheta/xqb/types"
)

// With adds a Common Table Expression (CTE)
func (qb *QueryBuilder) With(name string, query *QueryBuilder) *QueryBuilder {
	qb.withCTEs = append(qb.withCTEs, types.CTE{
		Name:  name,
		Query: query.GetData(),
	})
	return qb
}

// WithExpression adds a Common Table Expression using a raw SQL expression
func (qb *QueryBuilder) WithExpression(name string, sql string, bindings ...interface{}) *QueryBuilder {
	qb.withCTEs = append(qb.withCTEs, types.CTE{
		Name: name,
		Expression: &types.Expression{
			SQL:      sql,
			Bindings: bindings,
		},
	})
	return qb
}

// WithRecursive adds a recursive Common Table Expression (CTE)
func (qb *QueryBuilder) WithRecursive(name string, query *QueryBuilder) *QueryBuilder {
	qb.ctes = append(qb.ctes, types.CTE{
		Name:      name,
		Query:     query.GetData(),
		Recursive: true,
	})
	return qb
}

// WithRaw adds a raw Common Table Expression (CTE)
func (qb *QueryBuilder) WithRaw(name string, sql string, values ...interface{}) *QueryBuilder {
	rawQuery := &QueryBuilder{
		queryType: types.SELECT,
		columns:   []any{sql},
		bindings:  make([]types.Binding, len(values)),
	}
	for i, v := range values {
		rawQuery.bindings[i] = types.Binding{Value: v}
	}

	qb.ctes = append(qb.ctes, types.CTE{
		Name:  name,
		Query: rawQuery.GetData(),
	})
	return qb
}

// WithRecursiveRaw adds a raw recursive Common Table Expression (CTE)
func (qb *QueryBuilder) WithRecursiveRaw(name string, sql string, values ...interface{}) *QueryBuilder {
	rawQuery := &QueryBuilder{
		queryType: types.SELECT,
		columns:   []any{sql},
		bindings:  make([]types.Binding, len(values)),
	}
	for i, v := range values {
		rawQuery.bindings[i] = types.Binding{Value: v}
	}

	qb.ctes = append(qb.ctes, types.CTE{
		Name:      name,
		Query:     rawQuery.GetData(),
		Recursive: true,
	})
	return qb
}
