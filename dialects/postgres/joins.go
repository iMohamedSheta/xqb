package postgres

import (
	"github.com/iMohamedSheta/xqb/shared/types"
)

func (d *PostgresDialect) compileJoins(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql string

	for _, join := range qb.Joins {
		sql += " " + string(join.Type) + " " + d.Wrap(join.Table)

		if join.Type != types.CROSS_JOIN && join.Condition != "" {
			sql += " ON " + join.Condition
		}

		for _, binding := range join.Binding {
			bindings = append(bindings, binding.Value)
		}
	}

	return sql, bindings, nil
}
