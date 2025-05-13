package grammar

import (
	"strings"

	"github.com/iMohamedSheta/xqb/types"
)

// compileWhereClause compiles the WHERE clause
func (mg *MySQLGrammar) compileWhereClause(qb *types.QueryBuilderData) (string, []interface{}, error) {
	var bindings []interface{}
	var sql strings.Builder

	if len(qb.Where) > 0 {
		sql.WriteString(" WHERE ")
		for i, condition := range qb.Where {
			if i > 0 {
				sql.WriteString(" ")
				sql.WriteString(string(condition.Connector))
				sql.WriteString(" ")
			}

			if condition.Raw != nil {
				// Handle raw SQL expression
				sql.WriteString(condition.Raw.SQL)
				if len(condition.Raw.Bindings) > 0 {
					bindings = append(bindings, condition.Raw.Bindings...)
				}
			} else {
				// Handle regular condition
				sql.WriteString(condition.Column)
				if condition.Operator != "" {
					sql.WriteString(" ")
					sql.WriteString(condition.Operator)
					if condition.Value != nil {
						switch v := condition.Value.(type) {
						case []interface{}:
							// Handle IN and BETWEEN conditions
							if condition.Operator == "IN" || condition.Operator == "NOT IN" {
								placeholders := make([]string, len(v))
								for i := range v {
									placeholders[i] = "?"
									bindings = append(bindings, v[i])
								}
								sql.WriteString(" (")
								sql.WriteString(strings.Join(placeholders, ", "))
								sql.WriteString(")")
							} else if condition.Operator == "BETWEEN" || condition.Operator == "NOT BETWEEN" {
								if len(v) == 2 {
									sql.WriteString(" ? AND ?")
									bindings = append(bindings, v[0], v[1])
								}
							} else {
								sql.WriteString(" ?")
								bindings = append(bindings, v)
							}
						default:
							sql.WriteString(" ?")
							bindings = append(bindings, v)
						}
					}
				}
			}
		}
	}

	return sql.String(), bindings, nil
}
