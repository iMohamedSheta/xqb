package wrap

import (
	"fmt"
	"strconv"
	"strings"
)

func Wrap(value string, wrapChar byte) string {
	value = strings.TrimSpace(value)
	lower := strings.ToLower(value)

	// Handle aliases
	if idx := strings.LastIndex(lower, " as "); idx != -1 {
		left := strings.TrimSpace(value[:idx])
		right := strings.TrimSpace(value[idx+4:])
		return fmt.Sprintf("%s AS %s", Wrap(left, wrapChar), wrapValue(right, wrapChar))
	}

	// Handle shorthand aliases (e.g., users u)
	parts := strings.Fields(value)
	if len(parts) == 2 {
		return fmt.Sprintf("%s %s", wrapValue(parts[0], wrapChar), wrapValue(parts[1], wrapChar))
	}

	// Handle dot notation like table.column
	segments := strings.Split(value, ".")
	for i := range segments {
		segments[i] = wrapValue(segments[i], wrapChar)
	}
	return strings.Join(segments, ".")
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

func wrapValue(val string, wrapChar byte) string {
	val = strings.TrimSpace(val)

	if val == "*" {
		return "*"
	}

	wrap := string(wrapChar)

	if isWrapped(val, wrapChar) || isLiteral(val) || isLikelyExpr(val) {
		return val
	}

	// escape wraps in the value like my"value -> my""value
	escaped := strings.ReplaceAll(val, wrap, wrap+wrap)
	return wrap + escaped + wrap
}

func isWrapped(s string, wrapChar byte) bool {
	return strings.HasPrefix(s, string(wrapChar)) && strings.HasSuffix(s, string(wrapChar))
}
