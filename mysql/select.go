package mysql

import (
	"fmt"
	"strings"

	"github.com/iMohamedSheta/xqb/types"
)

// compileSelectClause compiles the SELECT clause
func (mg *MySQLGrammar) compileSelectClause(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql strings.Builder

	sql.WriteString("SELECT")

	if qb.IsHighPriority {
		sql.WriteString(" HIGH_PRIORITY")
	}

	if qb.IsStraightJoin {
		sql.WriteString(" STRAIGHT_JOIN")
	}

	if qb.IsCalcFoundRows {
		sql.WriteString(" SQL_CALC_FOUND_ROWS")
	}

	if qb.IsUsingDistinct {
		sql.WriteString(" DISTINCT")
	}

	// Handle columns
	if len(qb.Columns) == 0 && len(qb.AggregateFuncs) == 0 && len(qb.JSONExpressions) == 0 &&
		len(qb.StringFuncs) == 0 && len(qb.DateFuncs) == 0 && len(qb.MathExpressions) == 0 &&
		len(qb.ConditionalExprs) == 0 {
		sql.WriteString(" *")
	} else {
		columns := make([]string, 0)

		// Add regular columns
		for _, column := range qb.Columns {
			switch v := column.(type) {
			case string:
				columns = append(columns, v)
			case *types.Expression:
				columns = append(columns, v.SQL)
				bindings = append(bindings, v.Bindings...)
			default:
				columns = append(columns, fmt.Sprintf("%v", v))
			}
		}

		// Add aggregate functions
		for _, agg := range qb.AggregateFuncs {
			expr := string(agg.Function) + "("
			if agg.Distinct {
				expr += "DISTINCT "
			}
			expr += agg.Column + ")"
			if agg.Alias != "" {
				expr += " AS " + agg.Alias
			}
			columns = append(columns, expr)
		}

		// Add JSON expressions
		for _, json := range qb.JSONExpressions {
			funcName := "JSON_EXTRACT"
			if json.Function != "" {
				funcName = json.Function
			}
			expr := funcName + "(" + json.Column + ", '" + json.Path + "')"
			if json.Alias != "" {
				expr += " AS " + json.Alias
			}
			columns = append(columns, expr)
		}

		// Add string functions
		for _, str := range qb.StringFuncs {
			expr := str.Function + "(" + str.Column
			for _, param := range str.Params {
				expr += ", ?"
				bindings = append(bindings, param)
			}
			expr += ")"
			if str.Alias != "" {
				expr += " AS " + str.Alias
			}
			columns = append(columns, expr)
		}

		// Add date functions
		for _, date := range qb.DateFuncs {
			expr := date.Function + "(" + date.Column
			for _, param := range date.Params {
				expr += ", ?"
				bindings = append(bindings, param)
			}
			expr += ")"
			if date.Alias != "" {
				expr += " AS " + date.Alias
			}
			columns = append(columns, expr)
		}

		// Add math expressions
		for _, math := range qb.MathExpressions {
			expr := math.Expression
			if math.Alias != "" {
				expr += " AS " + math.Alias
			}
			columns = append(columns, expr)
		}

		// Add conditional expressions
		for _, cond := range qb.ConditionalExprs {
			expr := cond.Expression
			if cond.Alias != "" {
				expr += " AS " + cond.Alias
			}
			columns = append(columns, expr)
		}

		sql.WriteString(" ")
		sql.WriteString(strings.Join(columns, ", "))
	}

	return sql.String(), bindings, nil
}
