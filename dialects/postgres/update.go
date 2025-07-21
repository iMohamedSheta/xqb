package postgres

import (
	"errors"
	"fmt"
	"strings"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// CompileUpdate compiles the update operation for postgres driver
func (pg *PostgresDialect) CompileUpdate(qb *types.QueryBuilderData) (string, []any, error) {
	if len(qb.UpdatedBindings) == 0 {
		return "", nil, fmt.Errorf("%w: no bindings provided for update operation", xqbErr.ErrInvalidQuery)
	}

	var setParts []string

	var bindings []any
	var sql strings.Builder

	for _, binding := range qb.UpdatedBindings {
		setParts = append(setParts, fmt.Sprintf("%s = ?", binding.Column))
		bindings = append(bindings, binding.Value)
	}

	sql.WriteString(fmt.Sprintf("UPDATE %s SET %s", qb.Table.Name, strings.Join(setParts, ", ")))

	// Compile each part of the query in order
	clauses := []func(*types.QueryBuilderData) (string, []any, error){
		pg.compileWhereClause,
		pg.compileLimitClause,
	}

	for _, compiler := range clauses {
		if err := appendClause(&sql, &bindings, compiler, qb); err != nil {
			return "", nil, err
		}
	}

	// Check if there are any errors in building the query
	if len(qb.Errors) > 0 {
		errs := errors.Join(qb.Errors...)
		return "", nil, fmt.Errorf("%w: %s", xqbErr.ErrInvalidQuery, errs)
	}

	return sql.String(), bindings, nil
}
