package types

// QueryType represents the type of query being built
type QueryType int

// WhereConditionEnum represents the type of WHERE condition connector
type WhereConditionEnum string

// JoinType represents the type of JOIN
type JoinType string

// AggregateFunction represents an aggregate function
type AggregateFunction string

type UnionType string

const (
	SELECT QueryType = 1
	INSERT QueryType = 2
	UPDATE QueryType = 3
	DELETE QueryType = 4

	AND WhereConditionEnum = "AND"
	OR  WhereConditionEnum = "OR"

	INNER_JOIN JoinType = "INNER JOIN"
	LEFT_JOIN  JoinType = "LEFT JOIN"
	RIGHT_JOIN JoinType = "RIGHT JOIN"
	FULL_JOIN  JoinType = "FULL JOIN"
	CROSS_JOIN JoinType = "CROSS JOIN"

	SUM   AggregateFunction = "SUM"
	AVG   AggregateFunction = "AVG"
	MIN   AggregateFunction = "MIN"
	MAX   AggregateFunction = "MAX"
	COUNT AggregateFunction = "COUNT"

	UnionTypeUnion     UnionType = "Union"
	UnionTypeIntersect UnionType = "Intersect"
	UnionTypeExcept    UnionType = "Except"
)

// WhereCondition represents a WHERE clause condition
type WhereCondition struct {
	Column    string
	Operator  string
	Value     interface{}
	Connector WhereConditionEnum
	Raw       *Expression
	IsNot     bool
}

// OrderBy represents an ORDER BY clause
type OrderBy struct {
	Column    string
	Direction string
}

// Having represents a HAVING clause
type Having struct {
	Column    string
	Operator  string
	Value     interface{}
	Connector WhereConditionEnum
}

// Binding represents a value binding
type Binding struct {
	Column string
	Value  any
}

// Join represents a JOIN clause
type Join struct {
	Type      JoinType
	Table     string
	Condition string
	On        []WhereCondition
	Binding   []Binding
}

// Union represents a UNION clause
type Union struct {
	Expression *Expression
	Type       UnionType
	All        bool
}

// Expression represents a raw SQL expression
type Expression struct {
	SQL      string
	Bindings []interface{}
}

// CTE represents a Common Table Expression
type CTE struct {
	Name       string
	Query      interface{} // Will be *QueryBuilder
	Expression *Expression
	Recursive  bool
}

// AggregateExpr represents an aggregate function expression
type AggregateExpr struct {
	Function AggregateFunction
	Column   string
	Alias    string
	Distinct bool
}

// JSONExpression represents a JSON expression
type JSONExpression struct {
	Column   string
	Path     string
	Function string
	Alias    string
}

// MathExpression represents a mathematical expression
type MathExpression struct {
	Expression string
	Alias      string
}

// ConditionalExpr represents a conditional expression
type ConditionalExpr struct {
	Expression string
	Alias      string
}

// StringFunction represents a string function
type StringFunction struct {
	Function string
	Column   string
	Params   []interface{}
	Alias    string
}

// DateFunction represents a date function
type DateFunction struct {
	Function string
	Column   string
	Params   []interface{}
	Alias    string
}

// QueryBuilderData represents the data needed by grammars to compile queries
type QueryBuilderData struct {
	QueryType         QueryType
	Table             string
	Columns           []any
	ColumnAliases     map[string]string
	Where             []WhereCondition
	OrderBy           []OrderBy
	GroupBy           []string
	Having            []Having
	Limit             int
	Offset            int
	Joins             []Join
	Unions            []Union
	Bindings          []Binding
	Distinct          bool
	AggregateFuncs    []AggregateExpr
	Subqueries        map[string]interface{} // Will be *QueryBuilder
	WithCTEs          []CTE
	JSONExpressions   []JSONExpression
	MathExpressions   []MathExpression
	ConditionalExprs  []ConditionalExpr
	StringFuncs       []StringFunction
	DateFuncs         []DateFunction
	IndexHints        []string
	ForceIndex        string
	UseIndex          string
	IgnoreIndex       string
	IsUsingDistinct   bool
	IsForUpdate       bool
	IsLockInShareMode bool
	IsHighPriority    bool
	IsStraightJoin    bool
	IsCalcFoundRows   bool
	InsertedValues    []map[string]interface{} // Added for insert operations
}
