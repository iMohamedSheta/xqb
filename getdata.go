package xqb

import (
	"database/sql"
	"fmt"
	"math"
)

func (qb *QueryBuilder) Paginate(perPage int, page int, withCount bool) ([]map[string]any, map[string]any, error) {

	if page < 1 {
		page = 1
	}

	qb.limit = perPage
	qb.offset = (page - 1) * perPage

	results, err := qb.Execute(nil)

	meta := map[string]any{
		"per_page":     perPage,
		"current_page": page,
	}

	if withCount {
		// Use the Count() method to get the total count of records
		count := 5 // Need to add query that count the total

		if err != nil {
			return nil, nil, err
		}

		meta["total_count"] = count
		meta["last_page"] = int(math.Ceil(float64(count) / float64(perPage)))

		return results, meta, nil
	}

	// If no count is needed, just return the results
	return results, meta, nil
}

func (qb *QueryBuilder) Execute(tx *sql.Tx) ([]map[string]any, error) {

	qbData := qb.GetData()
	query, args, err := qb.grammar.Build(qbData)
	if err != nil {
		return nil, err
	}

	rows, err := executeQuery(tx, query, args)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	var results []map[string]any
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		result := make(map[string]any)
		for i, col := range columns {
			val := values[i]

			// Normalize value
			switch v := val.(type) {
			case []byte:
				result[col] = string(v)
			default:
				result[col] = v // leave everything else untouched
			}
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func executeQuery(tx *sql.Tx, query string, args ...any) (*sql.Rows, error) {
	// Flatten args if needed
	if len(args) == 1 {
		if a, ok := args[0].([]interface{}); ok {
			args = a
		}
	}

	if tx != nil {
		return tx.Query(query, args...)
	}

	dbManager := GetDBManager()

	if !dbManager.IsDBConnected() {
		return nil, ErrNoConnection
	}

	db, err := dbManager.GetDB()

	if err != nil {
		return nil, err
	}

	return db.Query(query, args...)
}

// Get executes the query and returns all results
func (qb *QueryBuilder) Get() ([]map[string]interface{}, error) {

	results, err := qb.Execute(nil)

	if err != nil {
		return nil, err
	}

	return results, nil
}

// First returns the first row
func (qb *QueryBuilder) First() (map[string]interface{}, error) {

	// Save current limit and reset after operation
	currentLimit := qb.limit
	defer func() { qb.limit = currentLimit }()

	qb.limit = 1

	results, err := qb.Get()
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, sql.ErrNoRows
	}

	return results[0], nil
}

// PluckMap gets a list of values for two columns where the first column becomes the key and the second becomes the value
func (qb *QueryBuilder) Pluck(value, key string) (map[string]any, error) {
	// Save current columns and reset after operation
	currentColumns := qb.columns
	defer func() { qb.columns = currentColumns }()

	qb.columns = []any{value, key}

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

// Value gets a single value from first row
func (qb *QueryBuilder) Value(column string) (interface{}, error) {
	// Save current columns and reset after operation
	currentColumns := qb.columns
	defer func() { qb.columns = currentColumns }()

	qb.columns = []any{column}

	result, err := qb.First()
	if err != nil {
		return nil, err
	}

	return result[column], nil
}
