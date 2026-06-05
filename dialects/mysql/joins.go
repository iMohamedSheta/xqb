package mysql

import (
	"fmt"
	"strings"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

func (d *MySqlDialect) compileJoins(qb *types.QueryBuilderData) (string, []any, error) {
	if len(qb.Joins) == 0 {
		return "", nil, nil
	}

	var bindings []any
	var parts []string

	for _, join := range qb.Joins {
		if join.Type == types.FULL_JOIN {
			return "", nil, fmt.Errorf("%w: FULL JOIN is not supported by MySQL dialect", xqbErr.ErrUnsupportedFeature)
		}

		part := string(join.Type) + " " + d.Wrap(join.Table)

		// Handle ON conditions (skip for CROSS JOIN)
		if join.Type != types.CROSS_JOIN {
			switch {
			case len(join.Conditions) > 0:
				// structured join conditions
				condSQL, condBindings, err := d.compileJoinConditions(join.Conditions, &bindings)
				if err != nil {
					return "", nil, err
				}
				part += " ON " + condSQL
				bindings = append(bindings, condBindings...)

			case join.Condition != "":
				// raw string condition
				part += " ON " + join.Condition
			}
		}

		// subquery/expression bindings
		for _, b := range join.Binding {
			bindings = append(bindings, b.Value)
		}

		parts = append(parts, part)
	}

	return " " + strings.Join(parts, " "), bindings, nil
}

func (d *MySqlDialect) compileJoinConditions(
	conditions []*types.JoinCondition,
	bindings *[]any,
) (string, []any, error) {

	var out []any
	var sb strings.Builder

	for i, cond := range conditions {
		if i > 0 {
			if cond.Connector == types.OR {
				sb.WriteString(" OR ")
			} else {
				sb.WriteString(" AND ")
			}
		}

		switch cond.Kind {

		case types.JoinConditionOn:
			sb.WriteString(fmt.Sprintf("%s %s %s",
				cond.First,
				cond.Operator,
				cond.Second,
			))

		case types.JoinConditionWhere:
			if cond.Value == nil {
				sb.WriteString(fmt.Sprintf("%s %s",
					cond.First,
					cond.Operator,
				))
			} else {
				sb.WriteString(fmt.Sprintf("%s %s ?",
					cond.First,
					cond.Operator,
				))
				out = append(out, cond.Value)
			}

		case types.JoinConditionRaw:
			sql, rawBindings, err := cond.Raw.ToSql("mysql")
			if err != nil {
				return "", nil, err
			}
			sb.WriteString(sql)
			out = append(out, rawBindings...)

		case types.JoinConditionGroup:
			groupSQL, groupBindings, err := d.compileJoinConditions(cond.Group, bindings)
			if err != nil {
				return "", nil, err
			}
			sb.WriteString("(" + groupSQL + ")")
			out = append(out, groupBindings...)

		default:
			return "", nil, fmt.Errorf("unsupported JoinConditionKind: %d", cond.Kind)
		}
	}

	return sb.String(), out, nil
}
