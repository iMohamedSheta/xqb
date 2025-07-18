package xqb

import (
	"fmt"
	"strings"

	"github.com/iMohamedSheta/xqb/shared/types"
)

type CaseWhen struct {
	cases      []string
	elseResult string
	alias      string
}

// Case creates a new CASE WHEN SQL expression builder.
func Case() *CaseWhen {
	return &CaseWhen{}
}

// When adds a WHEN condition to the CASE expression.
func (c *CaseWhen) When(condition string, result string) *CaseWhen {
	c.cases = append(c.cases, fmt.Sprintf("WHEN %s THEN %s", condition, result))
	return c
}

// Else adds an ELSE result to the CASE expression.
func (c *CaseWhen) Else(result string) *CaseWhen {
	c.elseResult = fmt.Sprintf("ELSE %s", result)
	return c
}

// As sets the alias for the CASE expression.
func (c *CaseWhen) As(alias string) *CaseWhen {
	c.alias = alias
	return c
}

// End builds the final CASE WHEN SQL expression.
func (c *CaseWhen) End() *types.Expression {
	raw := "CASE " + strings.Join(c.cases, " ") + " "
	if c.elseResult != "" {
		raw += c.elseResult + " "
	}
	raw += "END"
	if c.alias != "" {
		raw += " AS " + c.alias
	}
	return Raw(raw)
}
