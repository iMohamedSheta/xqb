package sqlserver

import (
	"fmt"
	"strconv"
	"strings"
)

// wrapBracket quotes identifiers with T-SQL [brackets]
func wrapBracket(value string) string {
	value = strings.TrimSpace(value)
	lower := strings.ToLower(value)

	// Handle aliases (users AS u)
	if idx := strings.LastIndex(lower, " as "); idx != -1 {
		left := strings.TrimSpace(value[:idx])
		right := strings.TrimSpace(value[idx+4:])
		return fmt.Sprintf("%s AS %s", wrapBracket(left), wrapBracketValue(right))
	}

	// Handle shorthand aliases (e.g., users u)
	parts := strings.Fields(value)
	if len(parts) == 2 {
		return fmt.Sprintf("%s %s", wrapBracketValue(parts[0], true), wrapBracketValue(parts[1], true))
	}

	// Handle dot notation like table.column / db.schema.table.column
	segments := strings.Split(value, ".")
	for i := range segments {
		segments[i] = wrapBracketValue(segments[i], false)
	}
	return strings.Join(segments, ".")
}

func wrapBracketValue(val string, _ ...bool) string {
	val = strings.TrimSpace(val)

	if val == "*" {
		return "*"
	}

	// table.* -> [table].*
	if strings.HasSuffix(val, ".*") {
		inner := strings.TrimSuffix(val, ".*")
		if inner == "" || inner == "*" {
			return val
		}
		return wrapBracketValue(inner) + ".*"
	}

	if isWrappedBracket(val) || isLiteral(val) || isLikelyExpr(val) {
		return val
	}

	escaped := strings.ReplaceAll(val, "]", "]]")
	return "[" + escaped + "]"
}

func isWrappedBracket(s string) bool {
	return strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]")
}

func isLiteral(s string) bool {
	if s == "null" || s == "true" || s == "false" {
		return true
	}
	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return true // numeric literal
	}
	if strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'") {
		return true // string literal
	}
	return false
}

func isLikelyExpr(s string) bool {
	return strings.ContainsAny(s, "()+*/-")
}
