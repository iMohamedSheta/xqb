package types

// Expression represents a raw Sql expression
type Expression struct {
	Sql      string
	Bindings []any
}

// DialectExpression represents a dialect expression
type DialectExpression struct {
	Default  string
	Dialects map[string]*Expression // dialect => expression
}

func (e DialectExpression) ToSql(dialect string) (string, []any, error) {
	if exp, ok := e.Dialects[dialect]; ok {
		return exp.ToSql()
	}
	return e.Dialects[e.Default].ToSql()
}

func (expr *Expression) ToSql() (string, []any, error) {
	return expr.Sql, expr.Bindings, nil
}
