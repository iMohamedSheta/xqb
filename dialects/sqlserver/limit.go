package sqlserver

import (
	"strconv"

	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileLimitClause compiles the FETCH NEXT clause.
// LIMIT without OFFSET is handled via SELECT TOP, so this only emits
// FETCH NEXT when both are present.
func (d *SqlServerDialect) compileLimitClause(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql string
	if qb.Limit != 0 && qb.Offset > 0 {
		sql += " FETCH NEXT " + strconv.Itoa(qb.Limit) + " ROWS ONLY"
	}
	return sql, bindings, nil
}

// compileOffsetClause compiles the OFFSET clause.
func (d *SqlServerDialect) compileOffsetClause(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql string
	if qb.Offset != 0 {
		sql += " OFFSET " + strconv.Itoa(qb.Offset) + " ROWS"
	}
	return sql, bindings, nil
}
