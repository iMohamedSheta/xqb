package grammar

import (
	"strconv"
	"strings"

	"github.com/iMohamedSheta/xqb/types"
)

// compileLimitOffsetClause compiles the LIMIT and OFFSET clauses
func (mg *MySQLGrammar) compileLimitOffsetClause(qb *types.QueryBuilderData) (string, []interface{}, error) {
	var bindings []interface{}
	var sql strings.Builder

	if qb.Limit > 0 || qb.Offset > 0 {
		if qb.Limit > 0 {
			sql.WriteString(" LIMIT ")
			sql.WriteString(strconv.Itoa(qb.Limit))
		}
		if qb.Offset > 0 {
			sql.WriteString(" OFFSET ")
			sql.WriteString(strconv.Itoa(qb.Offset))
		}
	}

	return sql.String(), bindings, nil
}
