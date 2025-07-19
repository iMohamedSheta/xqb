package types

// OrderBy represents an ORDER BY clause
type OrderBy struct {
	Column    string
	Direction string
	Raw       *Expression
}
