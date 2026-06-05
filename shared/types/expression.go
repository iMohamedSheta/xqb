package types

// types/expression.go

type ExpressionInterface interface {
	ToSql(dialect string) (sql string, bindings []any, err error)
}

// Expression — dialect-agnostic raw SQL
type Expression struct {
	Sql      string
	Bindings []any
}

// Make Expression satisfy ExpressionInterface (dialect param ignored)
func (e *Expression) ToSql(_ string) (string, []any, error) {
	return e.Sql, e.Bindings, nil
}

// DialectExpression — per-dialect SQL
type DialectExpression struct {
	Default  string
	Dialects map[string]*Expression
}

func (e *DialectExpression) ToSql(dialect string) (string, []any, error) {
	if exp, ok := e.Dialects[dialect]; ok {
		return exp.ToSql(dialect)
	}
	return e.Dialects[e.Default].ToSql(dialect)
}
