package mysql

import (
	"fmt"

	"github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileUnionClause compiles the union clauses for the postgres driver.
func (pg *MySqlDialect) compileUnionClause(qbd *types.QueryBuilderData) (string, []any, error) {
	var sql string
	var bindings []any
	// Add each union
	for _, union := range qbd.Unions {
		switch union.Type {
		case types.UnionTypeUnion:
			sql += " UNION "
		case types.UnionTypeIntersect, types.UnionTypeExcept:
			return "", nil, fmt.Errorf("%w: union type %s is not supported in MySql", errors.ErrUnsupportedFeature, string(union.Type))
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
