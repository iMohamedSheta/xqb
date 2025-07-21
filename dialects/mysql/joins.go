package mysql

import (
	"fmt"
	"strings"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

func (mg *MySQLGrammar) compileJoins(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql strings.Builder

	for _, join := range qb.Joins {
		if join.Type == types.FULL_JOIN {
			return "", nil, fmt.Errorf("%w: FULL JOIN is not supported by MySQL", xqbErr.ErrInvalidQuery)
		}

		sql.WriteString(" ")
		sql.WriteString(string(join.Type))
		sql.WriteString(" ")
		sql.WriteString(join.Table)

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
