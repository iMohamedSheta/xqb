package types

import (
	"github.com/iMohamedSheta/xqb/shared/enums"
)

// QueryBuilderData represents the data needed by grammars to compile queries
type QueryBuilderData struct {
	QueryType         enums.QueryType
	Table             *Table
	Columns           []any
	Where             []*WhereCondition
	OrderBy           []*OrderBy
	GroupBy           []string
	Having            []*Having
	Limit             int
	Offset            int
	Joins             []*Join
	Unions            []*Union
	Bindings          []*Binding
	Distinct          bool
	WithCTEs          []*CTE
	IsUsingDistinct   bool
	IsLockedForUpdate bool
	IsInSharedLock    bool
	InsertedValues    []map[string]any // Added for insert operations
	UpdatedBindings   []*Binding       // Added for update operations
	Errors            []error
}
