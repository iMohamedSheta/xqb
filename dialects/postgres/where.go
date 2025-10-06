package postgres

import (
	"fmt"
	"reflect"
	"strings"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

func (d *PostgresDialect) compileWhereClause(qb *types.QueryBuilderData) (string, []any, error) {
	if len(qb.Where) == 0 {
		return "", nil, nil
	}

	var sql string
	var bindings []any

	sql = " WHERE "
	for i, condition := range qb.Where {
		if i > 0 {
			sql += " " + string(condition.Connector) + " "
		}

		clause, b, err := d.compileWhereCondition(condition)
		if err != nil {
			return "", nil, err
		}
		sql += clause
		bindings = append(bindings, b...)
	}

	return sql, bindings, nil
}

func (d *PostgresDialect) compileWhereCondition(condition *types.WhereCondition) (string, []any, error) {
	if condition.Raw != nil {
		return condition.Raw.Sql, condition.Raw.Bindings, nil
	} else if len(condition.Group) > 0 {
		return d.compileGroupCondition(condition.Group, condition.Connector)
	}
	return d.compileBasicCondition(condition)
}

func (d *PostgresDialect) compileGroupCondition(group []*types.WhereCondition, connector types.WhereConditionEnum) (string, []any, error) {
	var sql string
	var bindings []any

	sql += "("
	for i, cond := range group {
		if i > 0 {
			sql += " " + string(cond.Connector) + " "
		}
		clause, b, err := d.compileWhereCondition(cond)
		if err != nil {
			return "", nil, err
		}
		sql += clause
		bindings = append(bindings, b...)
	}
	sql += ")"

	return sql, bindings, nil
}

func (d *PostgresDialect) compileBasicCondition(condition *types.WhereCondition) (string, []any, error) {
	var sql string
	var bindings []any

	sql += d.Wrap(condition.Column)

	if condition.Operator == "" {
		return "", nil, fmt.Errorf("%w: missing operator for column %q", xqbErr.ErrInvalidQuery, condition.Column)
	}

	op := strings.ToUpper(condition.Operator)
	sql += " " + op

	if condition.Value == nil {
		return sql, bindings, nil
	}

	switch op {
	case "IN", "NOT IN":
		rv := reflect.ValueOf(condition.Value)
		if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
			return "", nil, fmt.Errorf("%w: IN operator requires slice/array, got %T", xqbErr.ErrInvalidQuery, condition.Value)
		}

		n := rv.Len()
		if n == 0 {
			return "", nil, fmt.Errorf("%w: IN operator requires at least one value", xqbErr.ErrInvalidQuery)
		}

		placeholders := make([]string, n)
		for i := 0; i < n; i++ {
			bindings = append(bindings, rv.Index(i).Interface())
			placeholders[i] = "?"
		}
		sql += " (" + strings.Join(placeholders, ", ") + ")"

	case "BETWEEN", "NOT BETWEEN":
		rv := reflect.ValueOf(condition.Value)
		if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
			return "", nil, fmt.Errorf("%w: BETWEEN operator requires slice/array, got %T", xqbErr.ErrInvalidQuery, condition.Value)
		}
		if rv.Len() != 2 {
			return "", nil, fmt.Errorf("%w: BETWEEN operator requires exactly 2 values", xqbErr.ErrInvalidQuery)
		}

		sql += " ? AND ?"
		bindings = append(bindings, rv.Index(0).Interface(), rv.Index(1).Interface())

	default:
		sql += " ?"
		bindings = append(bindings, condition.Value)
	}

	return sql, bindings, nil
}
