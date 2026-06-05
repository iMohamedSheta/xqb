package postgres

import (
	"fmt"
	"strings"

	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileJoins compiles all JOIN clauses into SQL and returns the sql string and bindings
func (d *PostgresDialect) compileJoins(qb *types.QueryBuilderData) (string, []any, error) {
	if len(qb.Joins) == 0 {
		return "", nil, nil
	}

	var bindings []any
	var parts []string

	for _, join := range qb.Joins {
		part := string(join.Type) + " " + d.Wrap(join.Table)

		if join.Type != types.CROSS_JOIN {
			switch {
			case len(join.Conditions) > 0:
				// closure-based JoinClause: compile structured ON conditions
				condSQL, condBindings, err := d.compileJoinConditions(join.Conditions)
				if err != nil {
					return "", nil, err
				}
				part += " ON " + condSQL
				bindings = append(bindings, condBindings...)

			case join.Condition != "":
				// plain string condition: Join("orders", "users.id = orders.user_id")
				part += " ON " + join.Condition
			}
		}

		// table-level bindings come from subquery or expression joins
		for _, b := range join.Binding {
			bindings = append(bindings, b.Value)
		}

		parts = append(parts, part)
	}

	return " " + strings.Join(parts, " "), bindings, nil
}

// compileJoinConditions compiles a slice of JoinConditions into a SQL ON clause string and returns the sql string and bindings
func (d *PostgresDialect) compileJoinConditions(conditions []*types.JoinCondition) (string, []any, error) {
	var bindings []any
	var sb strings.Builder

	for i, cond := range conditions {
		// prepend connector for all conditions after the first
		if i > 0 {
			if cond.Connector == types.OR {
				sb.WriteString(" OR ")
			} else {
				sb.WriteString(" AND ")
			}
		}

		switch cond.Kind {
		case types.JoinConditionOn:
			// column op column — no binding: On("users.id", "=", "orders.user_id")
			sb.WriteString(fmt.Sprintf("%s %s %s", cond.First, cond.Operator, cond.Second))

		case types.JoinConditionWhere:
			if cond.Value == nil {
				// IS NULL / IS NOT NULL — no binding: WhereNull("orders.deleted_at")
				sb.WriteString(fmt.Sprintf("%s %s", cond.First, cond.Operator))
			} else {
				// column op $N — bound value: Where("orders.type", "=", "invoice")
				sb.WriteString(fmt.Sprintf("%s %s $%d", cond.First, cond.Operator, len(bindings)+1))
				bindings = append(bindings, cond.Value)
			}

		case types.JoinConditionRaw:
			// raw SQL fragment — rewrite ? placeholders to $N continuing from current offset
			sql, rawBindings, err := cond.Raw.ToSql("postgres")
			if err != nil {
				return "", nil, err
			}
			sb.WriteString(sql)
			bindings = append(bindings, rawBindings...)

		case types.JoinConditionGroup:
			// nested group — recursively compile and wrap in parentheses
			groupSQL, groupBindings, err := d.compileJoinConditions(cond.Group)
			if err != nil {
				return "", nil, err
			}
			sb.WriteString("(" + groupSQL + ")")
			bindings = append(bindings, groupBindings...)

		default:
			return "", nil, fmt.Errorf("unsupported JoinConditionKind: %d", cond.Kind)
		}
	}

	return sb.String(), bindings, nil
}
