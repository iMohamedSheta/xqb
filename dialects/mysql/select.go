package mysql

import (
	"fmt"
	"strings"

	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileSelectClause compiles the SELECT clause
func (mg *MySqlDialect) compileSelectClause(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql strings.Builder

	sql.WriteString("SELECT")

	if qb.IsUsingDistinct {
		sql.WriteString(" DISTINCT")
	}

	// Handle columns
	if len(qb.Columns) == 0 {
		sql.WriteString(" *")
	} else {
		columns := make([]string, 0)

		// Add regular columns
		for _, column := range qb.Columns {
			switch v := column.(type) {
			case string:
				columns = append(columns, mg.Wrap(v))
			case *types.Expression:
				columns = append(columns, v.Sql)
				bindings = append(bindings, v.Bindings...)
			case *types.DialectExpression:
				sqlStr, sqlBindings, err := v.ToSql(mg.GetDriver().String())
				if err != nil {
					return "", nil, err
				}
				columns = append(columns, sqlStr)
				bindings = append(bindings, sqlBindings...)
			default:
				columns = append(columns, fmt.Sprintf("%v", v))
			}
		}

		sql.WriteString(" ")
		sql.WriteString(strings.Join(columns, ", "))
	}

	return sql.String(), bindings, nil
}
