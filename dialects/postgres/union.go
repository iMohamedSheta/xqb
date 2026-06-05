package postgres

import (
	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileUnionClause compiles the union clauses for the postgres dialect.
func (d *PostgresDialect) compileUnionClause(qbd *types.QueryBuilderData) (string, []any, error) {
	var sql string
	var bindings []any
	// Add each union
	for _, union := range qbd.Unions {
		switch union.Type {
		case types.UnionTypeUnion:
			sql += " UNION "
		case types.UnionTypeIntersect:
			sql += " INTERSECT "
		case types.UnionTypeExcept:
			sql += " EXCEPT "
		}

		if union.All {
			sql += "ALL "
		}

		// Add the union query
		exprSql, exprBinding, err := union.Expression.ToSql(d.Getdialect().String())
		if err != nil {
			return "", nil, err
		}
		// Add the union query
		sql += "("
		sql += exprSql
		sql += ")"

		if len(exprBinding) > 0 {
			bindings = append(bindings, exprBinding...)
		}
	}

	return sql, bindings, nil
}
