package grammar

import (
	"strconv"
	"strings"

	"github.com/iMohamedSheta/xqb/types"
)

// compileLimitClause compiles the LIMIT clause
func (mg *MySQLGrammar) compileLimitClause(qb *types.QueryBuilderData) (string, []interface{}, error) {
	var bindings []interface{}
	var sql strings.Builder
	if qb.Limit != 0 {
		sql.WriteString(" LIMIT ")
		sql.WriteString(strconv.Itoa(qb.Limit))
	}
	return sql.String(), bindings, nil
}

// compileOffsetClause compiles the OFFSET clause
func (mg *MySQLGrammar) compileOffsetClause(qb *types.QueryBuilderData) (string, []interface{}, error) {
	var bindings []interface{}
	var sql strings.Builder
	if qb.Offset != 0 {
		sql.WriteString(" OFFSET ")
		sql.WriteString(strconv.Itoa(qb.Offset))
	}
	return sql.String(), bindings, nil
}
