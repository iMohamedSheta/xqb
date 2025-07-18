package types

type UnionType string

const (
	UnionTypeUnion     UnionType = "Union"
	UnionTypeIntersect UnionType = "Intersect"
	UnionTypeExcept    UnionType = "Except"
)

// Union represents a UNION clause
type Union struct {
	Expression *Expression
	Type       UnionType
	All        bool
}
