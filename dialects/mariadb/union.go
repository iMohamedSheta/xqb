package mariadb

import (
	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileUnionClause compiles the union clauses for the mariadb dialect.
// MariaDB supports UNION, INTERSECT and EXCEPT (10.3+).
func (d *MariaDBDialect) compileUnionClause(qbd *types.QueryBuilderData) (string, []any, error) {
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
