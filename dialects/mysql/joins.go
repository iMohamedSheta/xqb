package mysql

import (
	"fmt"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

func (d *MySqlDialect) compileJoins(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql string

	for _, join := range qb.Joins {
		if join.Type == types.FULL_JOIN {
			return "", nil, fmt.Errorf("%w: FULL JOIN is not supported by MySql dialect", xqbErr.ErrUnsupportedFeature)
		}

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
