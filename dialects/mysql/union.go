package mysql

import (
	"fmt"

	"github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileUnionClause compiles the union clauses for the postgres dialect.
func (d *MySqlDialect) compileUnionClause(qbd *types.QueryBuilderData) (string, []any, error) {
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
