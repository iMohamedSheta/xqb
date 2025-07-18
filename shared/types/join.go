package types

// JoinType represents the type of JOIN
type JoinType string

const (
	INNER_JOIN JoinType = "JOIN"
	LEFT_JOIN  JoinType = "LEFT JOIN"
	RIGHT_JOIN JoinType = "RIGHT JOIN"
	FULL_JOIN  JoinType = "FULL JOIN"
	CROSS_JOIN JoinType = "CROSS JOIN"
)

// Join represents a JOIN clause
type Join struct {
	Type      JoinType
	Table     string
	Condition string
	Binding   []Binding
}
