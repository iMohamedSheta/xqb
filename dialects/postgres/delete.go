package postgres

import (
	"errors"
	"fmt"
	"strings"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// CompileDelete compiles the delete operation for postgres driver
func (pg *PostgresDialect) CompileDelete(qb *types.QueryBuilderData) (string, []any, error) {
	tableName, _, err := pg.resolveTable(qb, "delete", false)
	if err != nil {
		return "", nil, err
	}

	// validate query builder delete build
	if err := pg.validateDelete(qb); err != nil {
		return "", nil, err
	}

	var bindings []any
	var sql strings.Builder

	// Compile the cte first
	cteSql, cteBindings, err := pg.compileCTEs(qb)
	if err != nil {
		return "", nil, err
	}
	sql.WriteString(cteSql)
	if cteBindings != nil {
		bindings = append(bindings, cteBindings...)
	}

	sql.WriteString(fmt.Sprintf("DELETE FROM %s", tableName))

	// Compile each part of the query in order
	clauses := []func(*types.QueryBuilderData) (string, []any, error){
		pg.compileWhereClause,
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

// validateDelete validates the delete operation for postgres driver
func (pg *PostgresDialect) validateDelete(qb *types.QueryBuilderData) error {
	var errs []error

	if len(qb.Where) == 0 {
		errs = append(errs, errors.New("DELETE without WHERE clause is dangerous we don't allow that"))
	}

	if qb.Limit > 0 {
		errs = append(errs, errors.New("LIMIT is not allowed in DELETE in the Postgres driver"))
	}

	if len(qb.OrderBy) > 0 {
		errs = append(errs, errors.New("ORDER BY is not allowed in DELETE in the Postgres driver"))
	}

	if len(qb.Having) != 0 {
		errs = append(errs, errors.New("HAVING is not allowed in DELETE in the Postgres driver"))
	}

	if qb.Offset > 0 {
		errs = append(errs, errors.New("OFFSET is not allowed in DELETE in the Postgres driver"))
	}

	if len(qb.GroupBy) > 0 {
		errs = append(errs, errors.New("GROUP BY is not allowed in DELETE in the Postgres driver"))
	}

	if len(qb.Joins) > 0 {
		errs = append(errs, errors.New("JOINs (USING clause) not supported in DELETE Postgres driver"))
	}

	if len(qb.Unions) > 0 {
		errs = append(errs, errors.New("UNION is not allowed in DELETE in the Postgres driver"))
	}

	if len(qb.Columns) > 0 {
		errs = append(errs, errors.New("COLUMNS are not valid in DELETE queries"))
	}

	if qb.Distinct || qb.IsUsingDistinct {
		errs = append(errs, errors.New("DISTINCT is not valid in DELETE queries"))
	}

	if len(qb.Options) != 0 {
		for option := range qb.Options {
			errs = append(errs, errors.New(option.String()+" Options are Not supported in DELETE queries"))
		}
	}

	if len(errs) != 0 {
		return fmt.Errorf("%w: %s", xqbErr.ErrInvalidQuery, errors.Join(errs...))
	}

	return nil
}
