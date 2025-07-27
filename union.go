package xqb

import "github.com/iMohamedSheta/xqb/shared/types"

func (qb *QueryBuilder) Union(secondaryQuery ...*QueryBuilder) *QueryBuilder {
	qb.addUnions(types.UnionTypeUnion, false, secondaryQuery)
	return qb
}

// Union adds a UNION clause
func (qb *QueryBuilder) UnionRaw(sql string, bindings ...any) *QueryBuilder {
	qb.unions = append(qb.unions, &types.Union{
		Expression: Raw(sql, bindings...),
		All:        false,
		Type:       types.UnionTypeUnion,
	})
	return qb
}

// UnionAll adds a UNION ALL
func (qb *QueryBuilder) UnionAll(secondaryQuery ...*QueryBuilder) *QueryBuilder {
	qb.addUnions(types.UnionTypeUnion, true, secondaryQuery)
	return qb
}

// UnionAll adds a UNION ALL
func (qb *QueryBuilder) UnionAllRaw(sql string, bindings ...any) *QueryBuilder {
	qb.unions = append(qb.unions, &types.Union{
		Expression: Raw(sql, bindings...),
		All:        true,
		Type:       types.UnionTypeUnion,
	})
	return qb
}

func (qb *QueryBuilder) ExceptUnion(secondaryQuery ...*QueryBuilder) *QueryBuilder {
	qb.addUnions(types.UnionTypeExcept, false, secondaryQuery)
	return qb
}

func (qb *QueryBuilder) ExceptUnionAll(secondaryQuery ...*QueryBuilder) *QueryBuilder {
	qb.addUnions(types.UnionTypeExcept, true, secondaryQuery)
	return qb
}

// ExceptUnion adds a EXCEPT clause
func (qb *QueryBuilder) ExceptUnionRaw(sql string, all bool, bindings ...any) *QueryBuilder {
	qb.unions = append(qb.unions, &types.Union{
		Expression: Raw(sql, bindings...),
		All:        all,
		Type:       types.UnionTypeExcept,
	})
	return qb
}

func (qb *QueryBuilder) IntersectUnion(secondaryQuery ...*QueryBuilder) *QueryBuilder {
	qb.addUnions(types.UnionTypeIntersect, false, secondaryQuery)
	return qb
}

// IntersectUnionAll adds a INTERSECT all clause using a secondary query
func (qb *QueryBuilder) IntersectUnionAll(secondaryQuery ...*QueryBuilder) *QueryBuilder {
	qb.addUnions(types.UnionTypeIntersect, true, secondaryQuery)
	return qb
}

// IntersectUnionRaw adds a INTERSECT clause using a raw Sql
func (qb *QueryBuilder) IntersectUnionRaw(sql string, all bool, bindings ...any) *QueryBuilder {
	qb.unions = append(qb.unions, &types.Union{
		Expression: Raw(sql, bindings...),
		All:        all,
		Type:       types.UnionTypeIntersect,
	})
	return qb
}

func (qb *QueryBuilder) addUnions(unionType types.UnionType, all bool, secondaryQuery []*QueryBuilder) {
	for _, sub := range secondaryQuery {
		qb.unions = append(qb.unions, &types.Union{
			Expression: sub.SetDialect(qb.GetDialect().Getdialect()).ToRawExpr(),
			All:        all,
			Type:       unionType,
		})
	}
}
