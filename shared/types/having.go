package types

// Having represents a HAVING clause
type Having struct {
	Column    string
	Operator  string
	Value     any
	Connector WhereConditionEnum
	Raw       *Expression
}
