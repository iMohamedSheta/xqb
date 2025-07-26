package xqb

import (
	"strings"

	"github.com/iMohamedSheta/xqb/shared/types"
)

type CaseWhen struct {
	cases      []string
	bindings   []any
	elseResult string
	alias      string
}

// Case creates a new CASE WHEN Sql expression builder.
func Case() *CaseWhen {
	return &CaseWhen{}
}

// When adds a WHEN condition to the CASE expression.
func (c *CaseWhen) When(condition string, result any, bindings ...any) *CaseWhen {
	c.cases = append(c.cases, "WHEN "+condition+" THEN ?")
	c.bindings = append(c.bindings, bindings...)
	c.bindings = append(c.bindings, result)
	return c
}

// Else adds an ELSE result to the CASE expression.
func (c *CaseWhen) Else(result any) *CaseWhen {
	c.elseResult = "ELSE ?"
	c.bindings = append(c.bindings, result)
	return c
}

// As sets the alias for the CASE expression.
func (c *CaseWhen) As(alias string) *CaseWhen {
	c.alias = alias
	return c
}

// End builds the final CASE WHEN Sql expression.
func (c *CaseWhen) End() *types.Expression {
	var raw string
	if len(c.cases) == 0 {
		raw = "CASE "
	} else {
		raw = "CASE " + strings.Join(c.cases, " ") + " "
	}
	if c.elseResult != "" {
		raw += c.elseResult + " "
	}
	raw += "END"
	if c.alias != "" {
		raw += " AS " + c.alias
	}
	return Raw(raw, c.bindings...)
}
