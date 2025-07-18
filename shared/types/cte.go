package types

// CTE represents a Common Table Expression
type CTE struct {
	Name       string
	Query      any // Will be *QueryBuilder
	Expression *Expression
	Recursive  bool
}
