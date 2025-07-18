package types

// Expression represents a raw SQL expression
type Expression struct {
	SQL      string
	Bindings []any
}

// DialectExpression represents a dialect expression
type DialectExpression struct {
	Default  string
	Dialects map[string]*Expression // dialect => expression
}

func (e DialectExpression) ToSQL(dialect string) (string, []any, error) {
	if exp, ok := e.Dialects[dialect]; ok {
		return exp.ToSQL()
	}
	return e.Dialects[e.Default].ToSQL()
}

func (expr *Expression) ToSQL() (string, []any, error) {
	return expr.SQL, expr.Bindings, nil
}
