package xqb

import (
	"fmt"
	"strings"

	"github.com/iMohamedSheta/xqb/shared/types"
)

// Count - returns a COUNT aggregate expression with optional alias.
func Count(column, alias string) *types.Expression {
	raw := fmt.Sprintf("COUNT(%s)", column)
	if alias != "" {
		raw += " AS " + alias
	}

	return Raw(raw)
}

// Sum - returns a SUM aggregate expression with optional alias.
func Sum(column string, alias string) *types.Expression {
	raw := fmt.Sprintf("SUM(%s)", column)
	if alias != "" {
		raw += " AS " + alias
	}

	return Raw(raw)
}

// Sum - returns a SUM aggregate expression with optional alias.
func Avg(column string, alias string) *types.Expression {
	raw := fmt.Sprintf("AVG(%s)", column)
	if alias != "" {
		raw += " AS " + alias
	}
	return Raw(raw)
}

// Min - returns a MIN aggregate expression with optional alias.
func Min(column string, alias string) *types.Expression {
	raw := fmt.Sprintf("MIN(%s)", column)
	if alias != "" {
		raw += " AS " + alias
	}
	return Raw(raw)
}

// Max  - returns a MAX aggregate expression with optional alias.
func Max(column string, alias string) *types.Expression {
	raw := fmt.Sprintf("MAX(%s)", column)
	if alias != "" {
		raw += " AS " + alias
	}
	return Raw(raw)
}

// JsonExtract - builds a JSON_EXTRACT expression with a given path and alias.
func JsonExtract(column string, path string, alias string) *types.DialectExpression {
	// Ensure MySQL path starts with "$."
	if !strings.HasPrefix(path, "$.") {
		path = "$." + strings.TrimPrefix(path, ".")
	}
	// Build MySQL-style access path
	mysql := Raw(fmt.Sprintf("JSON_EXTRACT(%s, '%s')", column, path))

	// Build PostgreSQL-style access path
	pgPath := strings.TrimPrefix(path, "$.")
	keys := strings.Split(pgPath, ".")
	pgExpr := column
	for _, key := range keys[:len(keys)-1] {
		pgExpr = fmt.Sprintf("%s->'%s'", pgExpr, key)
	}
	lastKey := keys[len(keys)-1]
	pgExpr = fmt.Sprintf("%s->>'%s'", pgExpr, lastKey)
	pg := Raw(pgExpr)

	if alias != "" {
		mysql = Raw(fmt.Sprintf("%s AS %s", mysql.SQL, alias))
		pg = Raw(fmt.Sprintf("%s AS %s", pg.SQL, alias))
	}

	return &types.DialectExpression{
		Default: "mysql",
		Dialects: map[string]*types.Expression{
			"mysql":    mysql,
			"postgres": pg,
		},
	}
}

// JSONFunc - builds a custom JSON function call with arguments and alias.
func JSONFunc(fnName string, args []string, alias string) *types.Expression {
	raw := fmt.Sprintf("%s(%s)", fnName, strings.Join(args, ", "))
	if alias != "" {
		raw += " AS " + alias
	}
	return Raw(raw)
}

// Math - returns a raw mathematical SQL expression with alias.
func Math(rawExpr string, alias string) *types.Expression {
	if alias != "" {
		rawExpr += " AS " + alias
	}
	return Raw(rawExpr)
}

// Date - formats a column as a DATE with optional alias.
func Date(column string, alias string) *types.Expression {
	raw := fmt.Sprintf("DATE(%s)", column)
	if alias != "" {
		raw += " AS " + alias
	}
	return Raw(raw)
}

// Date - formats a column as a DATE with optional alias.
func DateDiff(a, b, alias string) *types.Expression {
	raw := fmt.Sprintf("DATEDIFF(%s, %s)", a, b)
	if alias != "" {
		raw += " AS " + alias
	}
	return Raw(raw)
}

// DateAdd - returns a DATE_ADD expression using interval and unit.
func DateAdd(date, interval, unit, alias string) *types.Expression {
	raw := fmt.Sprintf("DATE_ADD(%s, INTERVAL %s %s)", date, interval, unit)
	if alias != "" {
		raw += " AS " + alias
	}
	return Raw(raw)
}

// DateSub - returns a DATE_SUB expression using interval and unit.
func DateSub(date, interval, unit, alias string) *types.Expression {
	raw := fmt.Sprintf("DATE_SUB(%s, INTERVAL %s %s)", date, interval, unit)
	if alias != "" {
		raw += " AS " + alias
	}
	return Raw(raw)
}

// DateFormat - returns a DATE_FORMAT expression with format and alias.
func DateFormat(column, format, alias string) *types.DialectExpression {
	mysqlExpr := Raw(fmt.Sprintf("DATE_FORMAT(%s, '%s')", column, format))
	pgExpr := Raw(fmt.Sprintf("TO_CHAR(%s, '%s')", column, format)) // PostgreSQL

	if alias != "" {
		mysqlExpr = Raw(fmt.Sprintf("%s AS %s", mysqlExpr.SQL, alias))
		pgExpr = Raw(fmt.Sprintf("%s AS %s", pgExpr.SQL, alias))
	}

	dialects := map[string]*types.Expression{
		"mysql":    mysqlExpr,
		"postgres": pgExpr,
	}

	dialectExpr := RawDialect("mysql", dialects)
	return &dialectExpr
}

// Coalesce - returns the first non-null value from the list.
func Coalesce(args []string, alias string) *types.Expression {
	raw := fmt.Sprintf("COALESCE(%s)", strings.Join(args, ", "))
	if alias != "" {
		raw += " AS " + alias
	}
	return Raw(raw)
}

// Concat - joins columns with optional bindings and alias.
func Concat(columns []string, alias string, bindings ...any) *types.Expression {
	raw := fmt.Sprintf("CONCAT(%s)", strings.Join(columns, ", "))
	if alias != "" {
		raw += " AS " + alias
	}

	return Raw(raw, bindings...)
}

// Lower - converts string column to lowercase
func Lower(column, alias string) *types.Expression {
	raw := fmt.Sprintf("LOWER(%s)", column)
	if alias != "" {
		raw += " AS " + alias
	}
	return Raw(raw)
}

// Upper - converts string column to uppercase
func Upper(column, alias string) *types.Expression {
	raw := fmt.Sprintf("UPPER(%s)", column)
	if alias != "" {
		raw += " AS " + alias
	}
	return Raw(raw)
}

// Length - returns the number of characters in a string
func Length(column, alias string) *types.Expression {
	raw := fmt.Sprintf("LENGTH(%s)", column)
	if alias != "" {
		raw += " AS " + alias
	}
	return Raw(raw)
}

// Trim - removes leading and trailing whitespace
func Trim(column, alias string) *types.Expression {
	raw := fmt.Sprintf("TRIM(%s)", column)
	if alias != "" {
		raw += " AS " + alias
	}
	return Raw(raw)
}

// Replace - replaces all occurrences of a substring
func Replace(column, from, to, alias string) *types.Expression {
	raw := fmt.Sprintf("REPLACE(%s, %s, %s)", column, from, to)
	if alias != "" {
		raw += " AS " + alias
	}
	return Raw(raw)
}

// Substring - extracts a substring (use MySQL-compatible syntax)
func Substring(column string, start, length int, alias string) *types.Expression {
	raw := fmt.Sprintf("SUBSTRING(%s, %d, %d)", column, start, length)
	if alias != "" {
		raw += " AS " + alias
	}
	return Raw(raw)
}
