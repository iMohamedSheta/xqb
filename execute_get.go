package xqb

import (
	"errors"
	"fmt"
	"math"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
)

// Get executes the query and returns all results
func (qb *QueryBuilder) Get() ([]map[string]any, error) {
	query, args, err := qb.GetSql()
	if err != nil {
		return nil, fmt.Errorf("%w: Get() Failed to build the sql query, %v", xqbErr.ErrInvalidQuery, err)
	}

	rows, err := Sql(query, args...).
		WithContext(qb.ctx).
		WithAfterExec(qb.settings.GetOnAfterQueryExecution()).
		Connection(qb.connection).
		WithTx(qb.tx).
		Query()

	if err != nil {
		return nil, fmt.Errorf("%w: Get() Invalid query sql query error %v", xqbErr.ErrQueryFailed, err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("%w: Get() failed to retrieve columns %v", xqbErr.ErrInvalidResult, err)
	}

	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	var results []map[string]any
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("%w: Get() failed to scan result rows %v", xqbErr.ErrInvalidResult, err)
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
		return nil, fmt.Errorf("%w: Get() failed to scan result rows %v", xqbErr.ErrInvalidResult, err)
	}

	return results, nil
}

// GetSql returns the sql query for Get()
func (qb *QueryBuilder) GetSql() (string, []any, error) {
	return qb.ToSql()
}

// First returns the first row
func (qb *QueryBuilder) First() (map[string]any, error) {
	qb.limit = 1

	results, err := qb.Get()
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, xqbErr.ErrNotFound
	}

	return results[0], nil
}

// FirstSql returns the sql query for First()
func (qb *QueryBuilder) FirstSql() (string, []any, error) {
	originalLimit := qb.limit
	qb.limit = 1
	defer func() { qb.limit = originalLimit }()

	return qb.ToSql()
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
		return nil, fmt.Errorf("%w: Value() column %q not found in result", xqbErr.ErrInvalidResult, column)
	}

	return val, nil
}

// ValueSql returns the sql query for Value()
func (qb *QueryBuilder) ValueSql(column string) (string, []any, error) {
	originalCols := qb.columns
	qb.columns = []any{column}
	defer func() { qb.columns = originalCols }()

	return qb.FirstSql()
}

// Paginate returns paginated results with optional count metadata
func (qb *QueryBuilder) Paginate(perPage int, page int, countBy string) ([]map[string]any, map[string]any, error) {
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

	if countBy != "" {
		count, err := qb.Count(countBy)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: Paginate() failed to get count of the records: %v", xqbErr.ErrInvalidQuery, err)
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
		return fmt.Errorf("%w: Chunks() chunk size must be greater than 0", xqbErr.ErrInvalidQuery)
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
			return fmt.Errorf("%w: Chunks() failed to process chunk %v", xqbErr.ErrUnsupportedFeature, err)
		}

		offset += chunkSize
	}

	return nil
}

// PaginateSql returns the sql query for Paginate()
func (qb *QueryBuilder) PaginateSql(perPage, page int) (string, []any, error) {
	if page < 1 {
		page = 1
	}

	qb.limit = perPage
	qb.offset = (page - 1) * perPage

	return qb.ToSql()
}

// Find finds the first result by ID
func (qb *QueryBuilder) Find(id any) (map[string]any, error) {
	return qb.Where("id", "=", id).First()
}

// FindSql returns the sql query for Find()
func (qb *QueryBuilder) FindSql(id any) (string, []any, error) {
	return qb.Where("id", "=", id).FirstSql()
}

// FindOrFail finds the first result by ID or returns a "not found" error
func (qb *QueryBuilder) FindOrFail(id any) (map[string]any, error) {
	result, err := qb.Find(id)
	if errors.Is(err, xqbErr.ErrNotFound) {
		return nil, fmt.Errorf("%w: FindOrFail() record with ID %v was not found", xqbErr.ErrNotFound, id)
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}
