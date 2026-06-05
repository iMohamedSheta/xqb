package types

// CTE represents a Common Table Expression
type CTE struct {
	Name       string
	Query      any // Will be *QueryBuilder
	Expression ExpressionInterface
	Recursive  bool
}
