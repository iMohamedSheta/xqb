package main

import "github.com/iMohamedSheta/xqb/types"

// Union adds a UNION clause
func (qb *QueryBuilder) Union(sql string, bindings ...interface{}) *QueryBuilder {
	qb.unions = append(qb.unions, types.Union{
		Expression: Raw(sql, bindings...),
		All:        false,
		Type:       types.UnionTypeUnion,
	})
	return qb
}

// UnionAll adds a UNION ALL
func (qb *QueryBuilder) UnionAll(sql string, bindings ...interface{}) *QueryBuilder {
	qb.unions = append(qb.unions, types.Union{
		Expression: Raw(sql, bindings...),
		All:        true,
		Type:       types.UnionTypeUnion,
	})
	return qb
}

// ExceptUnion adds a EXCEPT clause
func (qb *QueryBuilder) ExceptUnion(sql string, all bool, bindings ...interface{}) *QueryBuilder {
	qb.unions = append(qb.unions, types.Union{
		Expression: Raw(sql, bindings...),
		All:        all,
		Type:       types.UnionTypeExcept,
	})
	return qb
}

// IntersectUnion adds a INTERSECT clause
func (qb *QueryBuilder) IntersectUnion(sql string, all bool, bindings ...interface{}) *QueryBuilder {
	qb.unions = append(qb.unions, types.Union{
		Expression: Raw(sql, bindings...),
		All:        all,
		Type:       types.UnionTypeIntersect,
	})
	return qb
}
