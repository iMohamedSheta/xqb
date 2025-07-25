package xqb

import (
	"fmt"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
)

// PluckMap gets a list of values for two columns where the first column becomes the key and the second becomes the value
func (qb *QueryBuilder) Pluck(value, key string) (map[string]any, error) {

	err := qb.updateQueryForPluck(value, key)
	if err != nil {
		return nil, err
	}

	results, err := qb.Get()
	if err != nil {
		return nil, err
	}

	mappedResults := make(map[string]any)

	for _, row := range results {
		key, ok := row[key].(string)
		if !ok {
			// Try to convert the key to string
			keyStr := fmt.Sprintf("%v", row[key])
			mappedResults[keyStr] = row[value]
			continue
		}
		mappedResults[key] = row[value]
	}

	return mappedResults, nil
}

func (qb *QueryBuilder) PluckSQL(value, key string) (string, []any, error) {
	err := qb.updateQueryForPluck(value, key)
	if err != nil {
		return "", nil, err
	}

	sqlString, bindings, err := qb.ToSQL()
	if err != nil {
		return "", nil, err
	}

	return sqlString, bindings, nil
}

func (qb *QueryBuilder) updateQueryForPluck(value, key string) error {
	if value != "" && key != "" {
		qb.columns = []any{value, key}
	} else if value != "" {
		qb.columns = []any{value}
	} else if key != "" {
		qb.columns = []any{key}
	} else if len(qb.columns) == 0 {
		return fmt.Errorf("%w: Pluck() either value or key must be specified", xqbErr.ErrInvalidQuery)
	}

	return nil
}
