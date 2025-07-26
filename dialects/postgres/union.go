package postgres

import (
	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileUnionClause compiles the union clauses for the postgres driver.
func (pg *PostgresDialect) compileUnionClause(qbd *types.QueryBuilderData) (string, []any, error) {
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
		sql += "("
		sql += union.Expression.Sql
		sql += ")"

		if len(union.Expression.Bindings) > 0 {
			bindings = append(bindings, union.Expression.Bindings...)
		}
	}

	return sql, bindings, nil
}
