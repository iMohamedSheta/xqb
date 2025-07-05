package mysql

import (
	"strings"

	"github.com/iMohamedSheta/xqb/types"
)

// compileGroupByClause compiles the GROUP BY clause
func (mg *MySQLGrammar) compileGroupByClause(qb *types.QueryBuilderData) (string, []interface{}, error) {
	var bindings []interface{}
	var sql strings.Builder

	if len(qb.GroupBy) > 0 {
		sql.WriteString(" GROUP BY ")
		for i, column := range qb.GroupBy {
			if i > 0 {
				sql.WriteString(", ")
			}
			sql.WriteString(column)
		}
	}

	return sql.String(), bindings, nil
}

// compileHavingClause compiles the HAVING clause
func (mg *MySQLGrammar) compileHavingClause(qb *types.QueryBuilderData) (string, []interface{}, error) {
	var bindings []interface{}
	var sql strings.Builder

	if len(qb.Having) > 0 {
		sql.WriteString(" HAVING ")
		for i, having := range qb.Having {
			if i > 0 {
				sql.WriteString(" ")
				sql.WriteString(string(having.Connector))
				sql.WriteString(" ")
			}
			sql.WriteString(having.Column)
			sql.WriteString(" ")
			sql.WriteString(having.Operator)
			sql.WriteString(" ?")
			bindings = append(bindings, having.Value)
		}
	}

	return sql.String(), bindings, nil
}
