package xqb

import (
	"fmt"
)

// Count - Returns the number of rows in the result set.
func (qb *QueryBuilder) Count(column string) (int64, error) {
	qb.columns = []any{fmt.Sprintf("COUNT(%s) as count_value", column)}

	data, err := qb.Execute()
	if err != nil {
		return 0, err
	}

	if len(data) != 1 {
		return 0, fmt.Errorf("expected 1 row, got %d", len(data))
	}

	count, ok := data[0]["count_value"].(int64)
	if !ok {
		return 0, fmt.Errorf("expected count to be uint64, got %T", data[0]["count_value"])
	}

	return count, nil
}

// Avg - Returns the average value of a column.
func (qb *QueryBuilder) Avg(column string) (float64, error) {
	qb.columns = []any{fmt.Sprintf("AVG(%s) as avg_value", column)}

	data, err := qb.Execute()
	if err != nil {
		return 0, err
	}

	if len(data) != 1 {
		return 0, fmt.Errorf("expected 1 row, got %d", len(data))
	}

	avg, ok := data[0]["avg_value"].(float64)
	if !ok {
		return 0, fmt.Errorf("expected avg to be float64, got %T", data[0]["avg_value"])
	}

	return avg, nil
}

// Sum - Returns the sum of a column.
func (qb *QueryBuilder) Sum(column string) (float64, error) {
	qb.columns = []any{fmt.Sprintf("SUM(%s) as sum_value", column)}

	data, err := qb.Execute()
	if err != nil {
		return 0, err
	}

	if len(data) != 1 {
		return 0, fmt.Errorf("expected 1 row, got %d", len(data))
	}

	sum, ok := data[0]["sum_value"].(float64)
	if !ok {
		return 0, fmt.Errorf("expected sum to be float64, got %T", data[0]["sum_value"])
	}

	return sum, nil
}

// Min - Returns the minimum value of a column.
func (qb *QueryBuilder) Min(column string) (float64, error) {
	qb.columns = []any{fmt.Sprintf("MIN(%s) as min_value", column)}

	data, err := qb.Execute()
	if err != nil {
		return 0, err
	}

	if len(data) != 1 {
		return 0, fmt.Errorf("expected 1 row, got %d", len(data))
	}

	min, ok := data[0]["min_value"].(float64)
	if !ok {
		return 0, fmt.Errorf("expected min to be float64, got %T", data[0]["min_value"])
	}

	return min, nil
}

// Max - Returns the maximum value of a column.
func (qb *QueryBuilder) Max(column string) (float64, error) {
	qb.columns = []any{fmt.Sprintf("MAX(%s) as max_value", column)}

	data, err := qb.Execute()
	if err != nil {
		return 0, err
	}

	if len(data) != 1 {
		return 0, fmt.Errorf("expected 1 row, got %d", len(data))
	}

	max, ok := data[0]["max_value"].(float64)
	if !ok {
		return 0, fmt.Errorf("expected max to be float64, got %T", data[0]["max_value"])
	}

	return max, nil
}
