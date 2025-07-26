package postgres

import (
	"strconv"

	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileLimitClause compiles the LIMIT clause
func (pg *PostgresDialect) compileLimitClause(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql string
	if qb.Limit != 0 {
		sql += " LIMIT " + strconv.Itoa(qb.Limit)
	}
	return sql, bindings, nil
}

// compileOffsetClause compiles the OFFSET clause
func (pg *PostgresDialect) compileOffsetClause(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql string
	if qb.Offset != 0 {
		sql += " OFFSET " + strconv.Itoa(qb.Offset)
	}
	return sql, bindings, nil
}
