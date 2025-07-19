package mysql

import (
	"strings"

	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileFromClause compiles the FROM clause
func (mg *MySQLGrammar) compileFromClause(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql strings.Builder

	if qb.Table == "" {
		return "", nil, nil
	}

	sql.WriteString(" FROM ")
	sql.WriteString(qb.Table)

	return sql.String(), bindings, nil
}
