package mysql

import (
	"github.com/iMohamedSheta/xqb/shared/types"
)

func (d *MySqlDialect) compileHavingClause(qb *types.QueryBuilderData) (string, []any, error) {
	if len(qb.Having) == 0 {
		return "", nil, nil
	}

	var (
		bindings []any
		sql      = " HAVING "
	)

	for i, having := range qb.Having {
		if i > 0 {
			sql += " " + string(having.Connector) + " "
		}

		if having.Raw != nil {
			expSQL, expBindings, err := having.Raw.ToSql(d.Getdialect().String())
			if err != nil {
				return "", nil, err
			}

			sql += expSQL
			bindings = append(bindings, expBindings...)
			continue
		}

		sql += d.Wrap(having.Column) + " " + having.Operator

		if having.Value != nil {
			sql += " ?"
			bindings = append(bindings, having.Value)
		}
	}

	return sql, bindings, nil
}
