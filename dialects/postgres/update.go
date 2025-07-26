package postgres

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// CompileUpdate compiles the update operation for postgres driver
func (pg *PostgresDialect) CompileUpdate(qb *types.QueryBuilderData) (string, []any, error) {
	tableName, _, err := pg.resolveTable(qb, "update", false)
	if err != nil {
		return "", nil, err
	}

	// validate query builder update build
	if err := pg.validateUpdate(qb); err != nil {
		return "", nil, err
	}

	var setParts []string

	var bindings []any
	var sql strings.Builder

	// Sort updated bindings by column name
	sort.SliceStable(qb.UpdatedBindings, func(i, j int) bool {
		return qb.UpdatedBindings[i].Column < qb.UpdatedBindings[j].Column
	})

	for _, binding := range qb.UpdatedBindings {
		if expr, ok := binding.Value.(*types.Expression); ok {
			setParts = append(setParts, fmt.Sprintf("%s = %s", pg.Wrap(binding.Column), expr.Sql))
			bindings = append(bindings, expr.Bindings...)
		} else {
			setParts = append(setParts, fmt.Sprintf("%s = ?", pg.Wrap(binding.Column)))
			bindings = append(bindings, binding.Value)
		}
	}

	sql.WriteString(fmt.Sprintf("UPDATE %s", tableName))

	if err := appendClause(&sql, &bindings, pg.compileJoins, qb); err != nil {
		return "", nil, err
	}

	sql.WriteString(" SET ")
	sql.WriteString(strings.Join(setParts, ", "))

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

	return sql.String(), bindings, nil
}

// validateUpdate checks if the query builder is valid for the update operation
func (pg *PostgresDialect) validateUpdate(qb *types.QueryBuilderData) error {
	var errs []error

	if len(qb.Where) == 0 && !qb.AllowDangerous {
		errs = append(errs, errors.New("UPDATE without WHERE clause is dangerous we don't allow that you can add AllowDangerous to allow it"))
	}

	if len(qb.UpdatedBindings) == 0 {
		errs = append(errs, errors.New("no updated fields provided for update operation"))
	}

	if len(qb.Having) != 0 {
		errs = append(errs, errors.New("HAVING is not allowed in UPDATE operations in the Postgres driver"))
	}

	if qb.Offset > 0 {
		errs = append(errs, errors.New("OFFSET is not allowed in UPDATE in the Postgres driver"))
	}

	if len(qb.GroupBy) > 0 {
		errs = append(errs, errors.New("GROUP BY is not allowed in UPDATE in the Postgres driver"))
	}

	if len(qb.Unions) > 0 {
		errs = append(errs, errors.New("UNION is not allowed in UPDATE in the Postgres driver"))
	}

	if len(qb.Columns) > 0 {
		errs = append(errs, errors.New("SELECT COLUMNS are not valid in UPDATE queries"))
	}

	if qb.Distinct || qb.IsUsingDistinct {
		errs = append(errs, errors.New("DISTINCT is not valid in UPDATE queries"))
	}

	if len(qb.Options) != 0 {
		for option := range qb.Options {
			errs = append(errs, fmt.Errorf("option %s is not supported in UPDATE queries", option))
		}
	}

	if len(errs) != 0 {
		return fmt.Errorf("%w: %s", xqbErr.ErrInvalidQuery, errors.Join(errs...))
	}

	return nil
}
