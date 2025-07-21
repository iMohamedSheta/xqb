package postgres

import (
	"strconv"
	"strings"

	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileLimitClause compiles the LIMIT clause
func (pg *PostgresDialect) compileLimitClause(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql strings.Builder
	if qb.Limit != 0 {
		sql.WriteString(" LIMIT ")
		sql.WriteString(strconv.Itoa(qb.Limit))
	}
	return sql.String(), bindings, nil
}

// compileOffsetClause compiles the OFFSET clause
func (pg *PostgresDialect) compileOffsetClause(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql strings.Builder
	if qb.Offset != 0 {
		sql.WriteString(" OFFSET ")
		sql.WriteString(strconv.Itoa(qb.Offset))
	}
	return sql.String(), bindings, nil
}
