package xqb

import (
	"fmt"
	"strconv"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
)

// Count - Returns the number of rows in the result set.
func (qb *QueryBuilder) Count(column string) (int64, error) {
	qb.columns = []any{fmt.Sprintf("COUNT(%s) as count_value", column)}

	data, err := qb.Get()
	if err != nil {
		return 0, err
	}

	if len(data) != 1 {
		return 0, fmt.Errorf("%w: Count() expected one row as result, got %d rows", xqbErr.ErrInvalidResult, len(data))
	}

	count, err := asInt64(data[0]["count_value"])
	if err != nil {
		return 0, fmt.Errorf("%w: Count() failed to convert result value to int, %v", xqbErr.ErrInvalidResult, err)
	}

	return count, nil
}

// Avg - Returns the average value of a column.
func (qb *QueryBuilder) Avg(column string) (float64, error) {
	qb.columns = []any{fmt.Sprintf("AVG(%s) as avg_value", column)}

	data, err := qb.Get()
	if err != nil {
		return 0, err
	}

	if len(data) != 1 {
		return 0, fmt.Errorf("%w: Avg() expected one row as result, got %d rows", xqbErr.ErrInvalidResult, len(data))
	}

	avg, err := asFloat64(data[0]["avg_value"])
	if err != nil {
		return 0, fmt.Errorf("%w: Avg() failed to convert result value to float, %v", xqbErr.ErrInvalidResult, err)
	}

	return avg, nil
}

// Sum - Returns the sum of a column.
func (qb *QueryBuilder) Sum(column string) (float64, error) {
	qb.columns = []any{fmt.Sprintf("SUM(%s) as sum_value", column)}

	data, err := qb.Get()
	if err != nil {
		return 0, err
	}

	if len(data) != 1 {
		return 0, fmt.Errorf("%w: Sum() expected one row as result, got %d rows", xqbErr.ErrInvalidResult, len(data))
	}

	sum, err := asFloat64(data[0]["sum_value"])
	if err != nil {
		return 0, fmt.Errorf("%w: Sum() failed to convert result value to float, %v", xqbErr.ErrInvalidResult, err)
	}

	return sum, nil
}

// Min - Returns the minimum value of a column.
func (qb *QueryBuilder) Min(column string) (float64, error) {
	qb.columns = []any{fmt.Sprintf("MIN(%s) as min_value", column)}

	data, err := qb.Get()
	if err != nil {
		return 0, err
	}

	if len(data) != 1 {
		return 0, fmt.Errorf("%w: Min() expected one row as result, got %d rows", xqbErr.ErrInvalidResult, len(data))
	}

	min, err := asFloat64(data[0]["min_value"])
	if err != nil {
		return 0, fmt.Errorf("%w: Min() failed to convert result value to float, %v", xqbErr.ErrInvalidResult, err)
	}

	return min, nil
}

// Max - Returns the maximum value of a column.
func (qb *QueryBuilder) Max(column string) (float64, error) {
	qb.columns = []any{fmt.Sprintf("MAX(%s) as max_value", column)}

	data, err := qb.Get()
	if err != nil {
		return 0, err
	}

	if len(data) != 1 {
		return 0, fmt.Errorf("%w: Max() expected one row as result, got %d rows", xqbErr.ErrInvalidResult, len(data))
	}

	max, err := asFloat64(data[0]["max_value"])
	if err != nil {
		return 0, fmt.Errorf("%w: Max() failed to convert result value to float, %v", xqbErr.ErrInvalidResult, err)
	}

	return max, nil
}

// asInt64 - Converts a value to an int64
func asInt64(value any) (int64, error) {
	switch v := value.(type) {
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case []uint8:
		str := string(v)
		n, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid int64 value from []byte %q: %v", str, err)
		}
		return n, nil
	default:
		return 0, fmt.Errorf("unsupported type %T for numeric conversion", v)
	}
}

// asFloat64 - Converts a value to a float64
func asFloat64(value any) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case []uint8:
		str := string(v)
		n, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid float64 value from []byte %q: %v", str, err)
		}
		return n, nil
	default:
		return 0, fmt.Errorf("unsupported type %T for float64 conversion", v)
	}
}
