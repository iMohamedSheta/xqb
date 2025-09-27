package xqb

import (
	"fmt"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
)

// Count - Returns the number of rows in the result set.
func (qb *QueryBuilder) Count(column string) (int64, error) {
	qb.columns = []any{fmt.Sprintf("COUNT(%s) as count", qb.Wrap(column))}

	data, err := qb.Get()
	if err != nil {
		return 0, err
	}

	if len(data) != 1 {
		return 0, fmt.Errorf("%w: Count() expected one row as result, got %d rows", xqbErr.ErrInvalidResult, len(data))
	}

	count, ok := toInt64(data[0]["count"])
	if !ok {
		return 0, fmt.Errorf("%w: Count() failed to convert result value to int, %v", xqbErr.ErrInvalidResult, err)
	}

	return count, nil
}

// Avg - Returns the average value of a column.
func (qb *QueryBuilder) Avg(column string) (float64, error) {
	qb.columns = []any{fmt.Sprintf("AVG(%s) as avg", qb.Wrap(column))}

	data, err := qb.Get()
	if err != nil {
		return 0, err
	}

	if len(data) != 1 {
		return 0, fmt.Errorf("%w: Avg() expected one row as result, got %d rows", xqbErr.ErrInvalidResult, len(data))
	}

	avg, ok := toFloat64(data[0]["avg"])
	if !ok {
		return 0, fmt.Errorf("%w: Avg() failed to convert result value to float, %v", xqbErr.ErrInvalidResult, err)
	}

	return avg, nil
}

// Sum - Returns the sum of a column.
func (qb *QueryBuilder) Sum(column string) (float64, error) {
	qb.columns = []any{fmt.Sprintf("SUM(%s) as sum", qb.Wrap(column))}

	data, err := qb.Get()
	if err != nil {
		return 0, err
	}

	if len(data) != 1 {
		return 0, fmt.Errorf("%w: Sum() expected one row as result, got %d rows", xqbErr.ErrInvalidResult, len(data))
	}

	sum, ok := toFloat64(data[0]["sum"])
	if !ok {
		return 0, fmt.Errorf("%w: Sum() failed to convert result value to float, %v", xqbErr.ErrInvalidResult, err)
	}

	return sum, nil
}

// Min - Returns the minimum value of a column.
func (qb *QueryBuilder) Min(column string) (float64, error) {
	qb.columns = []any{fmt.Sprintf("MIN(%s) as min", qb.Wrap(column))}

	data, err := qb.Get()
	if err != nil {
		return 0, err
	}

	if len(data) != 1 {
		return 0, fmt.Errorf("%w: Min() expected one row as result, got %d rows", xqbErr.ErrInvalidResult, len(data))
	}

	min, ok := toFloat64(data[0]["min"])
	if !ok {
		return 0, fmt.Errorf("%w: Min() failed to convert result value to float, %v", xqbErr.ErrInvalidResult, err)
	}

	return min, nil
}

// Max - Returns the maximum value of a column.
func (qb *QueryBuilder) Max(column string) (float64, error) {
	qb.columns = []any{fmt.Sprintf("MAX(%s) as max", qb.Wrap(column))}

	data, err := qb.Get()
	if err != nil {
		return 0, err
	}

	if len(data) != 1 {
		return 0, fmt.Errorf("%w: Max() expected one row as result, got %d rows", xqbErr.ErrInvalidResult, len(data))
	}

	max, ok := toFloat64(data[0]["max"])
	if !ok {
		return 0, fmt.Errorf("%w: Max() failed to convert result value to float, %v", xqbErr.ErrInvalidResult, err)
	}

	return max, nil
}

// CountSql - Returns the Sql and bindings for the COUNT aggregate function.
func (qb *QueryBuilder) CountSql(column string) (string, []any, error) {
	qb.columns = []any{fmt.Sprintf("COUNT(%s) as count", qb.Wrap(column))}
	return qb.GetSql()
}

// AvgSql - Returns the Sql and bindings for the AVG aggregate function.
func (qb *QueryBuilder) AvgSql(column string) (string, []any, error) {
	qb.columns = []any{fmt.Sprintf("AVG(%s) as avg", qb.Wrap(column))}
	return qb.GetSql()
}

// SumSql - Returns the Sql and bindings for the SUM aggregate function.
func (qb *QueryBuilder) SumSql(column string) (string, []any, error) {
	qb.columns = []any{fmt.Sprintf("SUM(%s) as sum", qb.Wrap(column))}
	return qb.GetSql()
}

// MinSql - Returns the Sql and bindings for the MIN aggregate function.
func (qb *QueryBuilder) MinSql(column string) (string, []any, error) {
	qb.columns = []any{fmt.Sprintf("MIN(%s) as min", qb.Wrap(column))}
	return qb.GetSql()
}

// MaxSql - Returns the Sql and bindings for the MAX aggregate function.
func (qb *QueryBuilder) MaxSql(column string) (string, []any, error) {
	qb.columns = []any{fmt.Sprintf("MAX(%s) as max", qb.Wrap(column))}
	return qb.GetSql()
}
