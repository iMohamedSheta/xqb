package types

import (
	"github.com/iMohamedSheta/xqb/shared/enums"
)

// QueryBuilderData represents the data needed by grammars to compile queries
type QueryBuilderData struct {
	QueryType       enums.QueryType
	Table           *Table
	Columns         []any
	Where           []*WhereCondition
	OrderBy         []*OrderBy
	GroupBy         []string
	Having          []*Having
	Limit           int
	Offset          int
	Joins           []*Join
	Unions          []*Union
	Bindings        []*Binding
	Distinct        bool
	WithCTEs        []*CTE
	IsUsingDistinct bool
	InsertedValues  []map[string]any // Added for insert operations
	UpdatedBindings []*Binding       // Added for update operations
	Errors          []error
	DeleteFrom      []string
	Options         map[Option]any // field for flexible Sql extensions
	AllowDangerous  bool
}

func (qb *QueryBuilderData) SetOption(key Option, value any) {
	qb.Options[key] = value
}

func (qb *QueryBuilderData) GetOption(key Option) (any, bool) {
	val, ok := qb.Options[key]
	return val, ok
}
