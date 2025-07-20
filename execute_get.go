package xqb

import (
	"errors"
	"fmt"
	"math"
)

// Get executes the query and returns all results
func (qb *QueryBuilder) Get() ([]map[string]any, error) {
	qbData := qb.GetData()
	query, args, err := qb.grammar.Build(qbData)
	if err != nil {
		return nil, fmt.Errorf("%w [Get]: Failed to build the sql, %v", ErrInvalidQuery, err)
	}

	rows, err := Sql(query, args...).Connection(qb.connection).WithTx(qb.tx).Query()
	if err != nil {
		return nil, fmt.Errorf("%w [Get]: Invalid query sql query error %v", ErrInvalidExecutedQuerySyntax, err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("%w [Get]: %v", ErrInvalidExecutedQuerySyntax, err)
	}

	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	var results []map[string]any
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("%w [Get]: %v", ErrInvalidResultType, err)
		}

		result := make(map[string]any)
		for i, col := range columns {
			val := values[i]
			switch v := val.(type) {
			case []byte:
				result[col] = string(v)
			default:
				result[col] = v
			}
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w [Get]: %v", ErrInvalidResultType, err)
	}

	return results, nil
}

// First returns the first row
func (qb *QueryBuilder) First() (map[string]any, error) {
	qb.limit = 1

	results, err := qb.Get()
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, ErrNotFound
	}

	return results[0], nil
}

// Value gets a single value from the first row
func (qb *QueryBuilder) Value(column string) (any, error) {
	qb.columns = []any{column}

	result, err := qb.First()
	if err != nil {
		return nil, err
	}

	val, ok := result[column]
	if !ok {
		return nil, fmt.Errorf("%w [Value]: column %q not found in result", ErrInvalidResultType, column)
	}

	return val, nil
}

// Paginate returns paginated results with optional count metadata
func (qb *QueryBuilder) Paginate(perPage int, page int, withCount bool) ([]map[string]any, map[string]any, error) {
	if page < 1 {
		page = 1
	}

	qb.limit = perPage
	qb.offset = (page - 1) * perPage

	results, err := qb.Get()
	if err != nil {
		return nil, nil, err
	}

	meta := map[string]any{
		"per_page":     perPage,
		"current_page": page,
	}

	if withCount {
		count, err := qb.Count("*")
		if err != nil {
			return nil, nil, fmt.Errorf("%w [Paginate]: failed to count records: %v", ErrInvalidQuery, err)
		}

		lastPage := int(math.Ceil(float64(count) / float64(perPage)))

		var nextPage, prevPage any
		if page < lastPage {
			nextPage = page + 1
		}
		if page > 1 {
			prevPage = page - 1
		}

		meta["total_count"] = count
		meta["last_page"] = lastPage
		meta["next_page"] = nextPage
		meta["prev_page"] = prevPage
	}

	return results, meta, nil
}

// Chunks processes results in batch and calls the closure for each chunk
func (qb *QueryBuilder) Chunks(chunkSize int, closure func(results []map[string]any) error) error {
	if chunkSize <= 0 {
		return fmt.Errorf("%w [Chunks]: chunk size must be greater than 0", ErrInvalidQuery)
	}

	offset := 0
	qb.limit = chunkSize

	for {
		qb.offset = offset

		results, err := qb.Get()
		if err != nil {
			return err
		}

		if len(results) == 0 {
			break
		}

		if err := closure(results); err != nil {
			return fmt.Errorf("%w [Chunks]: %v", ErrUnsupportedFeature, err)
		}

		offset += chunkSize
	}

	return nil
}

// Find finds the first result by ID
func (qb *QueryBuilder) Find(id any) (map[string]any, error) {
	return qb.Where("id", "=", id).First()
}

// FindOrFail finds the first result by ID or returns a "not found" error
func (qb *QueryBuilder) FindOrFail(id any) (map[string]any, error) {
	result, err := qb.Find(id)
	if errors.Is(err, ErrNotFound) {
		return nil, fmt.Errorf("%w [FindOrFail]: record with ID %v not found", ErrNotFound, id)
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}
