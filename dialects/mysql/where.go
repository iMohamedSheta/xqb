package mysql

import (
	"fmt"
	"strings"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

func (mg *MySQLDialect) compileWhereClause(qb *types.QueryBuilderData) (string, []any, error) {
	if len(qb.Where) == 0 {
		return "", nil, nil
	}

	var sql strings.Builder
	var bindings []any

	sql.WriteString(" WHERE ")
	for i, condition := range qb.Where {
		if i > 0 {
			sql.WriteString(" ")
			sql.WriteString(string(condition.Connector))
			sql.WriteString(" ")
		}

		clause, b, err := mg.compileWhereCondition(condition)
		if err != nil {
			return "", nil, err
		}
		sql.WriteString(clause)
		bindings = append(bindings, b...)
	}

	return sql.String(), bindings, nil
}

func (mg *MySQLDialect) compileWhereCondition(condition *types.WhereCondition) (string, []any, error) {
	if condition.Raw != nil {
		return condition.Raw.SQL, condition.Raw.Bindings, nil
	} else if len(condition.Group) > 0 {
		return mg.compileGroupCondition(condition.Group, condition.Connector)
	}
	return mg.compileBasicCondition(condition)
}

func (mg *MySQLDialect) compileGroupCondition(group []*types.WhereCondition, connector types.WhereConditionEnum) (string, []any, error) {
	var sql strings.Builder
	var bindings []any

	sql.WriteString("(")
	for i, cond := range group {
		if i > 0 {
			sql.WriteString(" ")
			sql.WriteString(string(cond.Connector))
			sql.WriteString(" ")
		}
		clause, b, err := mg.compileWhereCondition(cond)
		if err != nil {
			return "", nil, err
		}
		sql.WriteString(clause)
		bindings = append(bindings, b...)
	}
	sql.WriteString(")")

	return sql.String(), bindings, nil
}

func (mg *MySQLDialect) compileBasicCondition(condition *types.WhereCondition) (string, []any, error) {
	var sql strings.Builder
	var bindings []any

	sql.WriteString(mg.Wrap(condition.Column))

	if condition.Operator == "" {
		return "", nil, fmt.Errorf("%w: missing operator for column %q", xqbErr.ErrInvalidQuery, condition.Column)
	}

	op := strings.ToUpper(condition.Operator)

	sql.WriteString(" ")
	sql.WriteString(op)

	if condition.Value == nil {
		return sql.String(), bindings, nil
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
			sql.WriteString(" (")
			sql.WriteString(strings.Join(placeholders, ", "))
			sql.WriteString(")")
		case "BETWEEN", "NOT BETWEEN":
			if len(v) != 2 {
				return "", nil, fmt.Errorf("%w: BETWEEN operator requires exactly 2 values", xqbErr.ErrInvalidQuery)
			}
			sql.WriteString(" ? AND ?")
			bindings = append(bindings, v[0], v[1])
		default:
			return "", nil, fmt.Errorf("%w: unsupported operator %q for slice value in column %q", xqbErr.ErrInvalidQuery, op, condition.Column)
		}
	default:
		sql.WriteString(" ?")
		bindings = append(bindings, v)
	}

	return sql.String(), bindings, nil
}
