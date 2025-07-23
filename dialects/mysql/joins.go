package mysql

import (
	"fmt"
	"strings"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

func (mg *MySQLDialect) compileJoins(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql strings.Builder

	for _, join := range qb.Joins {
		if join.Type == types.FULL_JOIN {
			return "", nil, fmt.Errorf("%w: FULL JOIN is not supported by MySQL driver", xqbErr.ErrUnsupportedFeature)
		}

		sql.WriteString(" ")
		sql.WriteString(string(join.Type))
		sql.WriteString(" ")
		sql.WriteString(mg.Wrap(join.Table))

		if join.Type != types.CROSS_JOIN && join.Condition != "" {
			sql.WriteString(" ON ")
			sql.WriteString(join.Condition)
		}

		for _, binding := range join.Binding {
			bindings = append(bindings, binding.Value)
		}
	}

	return sql.String(), bindings, nil
}
