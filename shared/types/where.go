package types

// WhereConditionEnum represents the type of WHERE condition connector
type WhereConditionEnum string

const (
	AND WhereConditionEnum = "AND"
	OR  WhereConditionEnum = "OR"
)

// WhereCondition represents a WHERE clause condition
type WhereCondition struct {
	Column    string
	Operator  string
	Value     any
	Connector WhereConditionEnum
	Raw       *Expression
}
