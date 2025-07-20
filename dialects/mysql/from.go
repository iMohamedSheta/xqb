package mysql

import (
	"fmt"
	"strings"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileFromClause compiles the FROM clause
func (mg *MySQLGrammar) compileFromClause(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql strings.Builder

	if qb.Table == nil || (qb.Table.Raw == nil && qb.Table.Name == "") {
		qb.Errors = append(qb.Errors, fmt.Errorf("%w: table name is required for select statement", xqbErr.ErrInvalidQuery))
		return "", nil, nil
	}

	if qb.Table.Raw != nil && qb.Table.Name != "" {
		qb.Errors = append(qb.Errors, fmt.Errorf("%w: both raw SQL and table name set; choose one", xqbErr.ErrInvalidQuery))
		return "", nil, nil
	}

	sql.WriteString(" FROM ")

	if qb.Table.Raw != nil {
		sql.WriteString(qb.Table.Raw.SQL)
		bindings = append(bindings, qb.Table.Raw.Bindings...)
	} else {
		sql.WriteString(qb.Table.Name)
	}

	return sql.String(), bindings, nil
}
