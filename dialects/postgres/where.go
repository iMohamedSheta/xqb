package postgres

import (
	"fmt"
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

	switch v := condition.Value.(type) {
	case []any:
		switch op {
		case "IN", "NOT IN":
			if len(v) == 0 {
				return "", nil, fmt.Errorf("%w: IN operator requires at least one value", xqbErr.ErrInvalidQuery)
			}
			placeholders := make([]string, len(v))
			for i := range v {
				placeholders[i] = "?"
				bindings = append(bindings, v[i])
			}
			sql += " (" + strings.Join(placeholders, ", ") + ")"
		case "BETWEEN", "NOT BETWEEN":
			if len(v) != 2 {
				return "", nil, fmt.Errorf("%w: BETWEEN operator requires exactly 2 values", xqbErr.ErrInvalidQuery)
			}
			sql += " ? AND ?"
			bindings = append(bindings, v[0], v[1])
		default:
			return "", nil, fmt.Errorf("%w: unsupported operator %q for slice value in column %q", xqbErr.ErrInvalidQuery, op, condition.Column)
		}
	default:
		sql += " ?"
		bindings = append(bindings, v)
	}

	return sql, bindings, nil
}
