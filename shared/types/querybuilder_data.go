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
	if qb.Options == nil {
		qb.Options = make(map[Option]any)
	}
	qb.Options[key] = value
}

func (qb *QueryBuilderData) GetOption(key Option) (any, bool) {
	if qb.Options == nil {
		return nil, false
	}
	val, ok := qb.Options[key]
	return val, ok
}

// Type-safe helper methods

func (qb *QueryBuilderData) GetBoolOption(key Option) (bool, bool) {
	val, ok := qb.GetOption(key)
	if !ok {
		return false, false
	}
	boolVal, isBool := val.(bool)
	return boolVal, isBool
}

func (qb *QueryBuilderData) GetStringSliceOption(key Option) ([]string, bool) {
	val, ok := qb.GetOption(key)
	if !ok {
		return nil, false
	}
	sliceVal, isSlice := val.([]string)
	return sliceVal, isSlice
}

func (qb *QueryBuilderData) GetLockOption(key Option) (OptionValueLock, bool) {
	val, ok := qb.GetOption(key)
	if !ok {
		return 0, false
	}
	lockVal, isLock := val.(OptionValueLock)
	return lockVal, isLock
}

func (qb *QueryBuilderData) GetLockWaitOption(key Option) (OptionValueLockWait, bool) {
	val, ok := qb.GetOption(key)
	if !ok {
		return 0, false
	}
	waitVal, isWait := val.(OptionValueLockWait)
	return waitVal, isWait
}
