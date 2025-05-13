package grammar

import (
	"strings"

	"github.com/iMohamedSheta/xqb/types"
)

// compileFromClause compiles the FROM clause
func (mg *MySQLGrammar) compileFromClause(qb *types.QueryBuilderData) (string, []interface{}, error) {
	var bindings []interface{}
	var sql strings.Builder

	if qb.Table == "" {
		return "", nil, nil
	}

	sql.WriteString(" FROM ")
	sql.WriteString(qb.Table)

	// Add index hints
	if qb.ForceIndex != "" {
		sql.WriteString(" FORCE INDEX (")
		sql.WriteString(qb.ForceIndex)
		sql.WriteString(")")
	} else if qb.UseIndex != "" {
		sql.WriteString(" USE INDEX (")
		sql.WriteString(qb.UseIndex)
		sql.WriteString(")")
	} else if qb.IgnoreIndex != "" {
		sql.WriteString(" IGNORE INDEX (")
		sql.WriteString(qb.IgnoreIndex)
		sql.WriteString(")")
	}

	return sql.String(), bindings, nil
}
