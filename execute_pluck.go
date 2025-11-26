package xqb

import (
	"fmt"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
)

// PluckSlice returns []any of a single column
func (qb *QueryBuilder) PluckSlice(valueField string) ([]any, error) {
	if valueField == "" {
		return nil, fmt.Errorf("%w: PluckSlice() requires a value field", xqbErr.ErrInvalidQuery)
	}

	if err := qb.updateQueryForPluck(valueField); err != nil {
		return nil, err
	}

	results, err := qb.Get()
	if err != nil {
		return nil, err
	}

	values := make([]any, len(results))
	for i, row := range results {
		values[i] = row[valueField]
	}

	return values, nil
}

// PluckMap gets a list of values for two columns where the first column becomes the key and the second becomes the value
func (qb *QueryBuilder) PluckMap(valueField, keyField string) (map[string]any, error) {
	if valueField == "" || keyField == "" {
		return nil, fmt.Errorf("%w: PluckMap() requires both value and key fields", xqbErr.ErrInvalidQuery)
	}

	if err := qb.updateQueryForPluck(valueField, keyField); err != nil {
		return nil, err
	}

	results, err := qb.Get()
	if err != nil {
		return nil, err
	}

	mappedResults := make(map[string]any)
	for _, row := range results {
		key := fmt.Sprintf("%v", row[keyField])
		mappedResults[key] = row[valueField]
	}

	return mappedResults, nil
}

func (qb *QueryBuilder) PluckSliceSql(valueField string) (string, []any, error) {
	if valueField == "" {
		return "", nil, fmt.Errorf("%w: PluckSlice() requires a value field", xqbErr.ErrInvalidQuery)
	}

	err := qb.updateQueryForPluck(valueField)
	if err != nil {
		return "", nil, err
	}

	sqlString, bindings, err := qb.ToSql()
	if err != nil {
		return "", nil, err
	}

	return sqlString, bindings, nil
}

func (qb *QueryBuilder) PluckMapSql(valueField, keyField string) (string, []any, error) {
	if valueField == "" || keyField == "" {
		return "", nil, fmt.Errorf("%w: PluckMap() requires both value and key fields", xqbErr.ErrInvalidQuery)
	}

	err := qb.updateQueryForPluck(valueField, keyField)
	if err != nil {
		return "", nil, err
	}

	sqlString, bindings, err := qb.ToSql()
	if err != nil {
		return "", nil, err
	}

	return sqlString, bindings, nil
}

func (qb *QueryBuilder) updateQueryForPluck(value string, key ...string) error {
	if value == "" {
		return fmt.Errorf("%w: Pluck() requires a value field", xqbErr.ErrInvalidQuery)
	}

	if len(key) > 0 {
		qb.columns = []any{value, key[0]}
	} else {
		qb.columns = []any{value}
	}

	return nil
}
