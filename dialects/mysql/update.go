package mysql

import (
	"errors"
	"fmt"
	"strings"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// CompileUpdate compiles the update query
func (mg *MySqlDialect) CompileUpdate(qb *types.QueryBuilderData) (string, []any, error) {
	tableName, _, err := mg.resolveTable(qb, "update", false)
	if err != nil {
		return "", nil, err
	}

	// validate query builder update build
	if err := mg.validateUpdate(qb); err != nil {
		return "", nil, err
	}

	var setParts []string

	var bindings []any
	var sql strings.Builder

	for _, binding := range qb.UpdatedBindings {
		setParts = append(setParts, fmt.Sprintf("%s = ?", mg.Wrap(binding.Column)))
		bindings = append(bindings, binding.Value)
	}

	sql.WriteString(fmt.Sprintf("UPDATE %s", tableName))

	if err := appendClause(&sql, &bindings, mg.compileJoins, qb); err != nil {
		return "", nil, err
	}

	sql.WriteString(" SET ")
	sql.WriteString(strings.Join(setParts, ", "))

	// Compile each part of the query in order
	clauses := []func(*types.QueryBuilderData) (string, []any, error){
		mg.compileWhereClause,
		mg.compileOrderByClause,
		mg.compileLimitClause,
	}

	for _, compiler := range clauses {
		if err := appendClause(&sql, &bindings, compiler, qb); err != nil {
			return "", nil, err
		}
	}

	return sql.String(), bindings, nil
}

// validateUpdate checks if the query builder is valid for the update operation
func (mg *MySqlDialect) validateUpdate(qb *types.QueryBuilderData) error {
	var errs []error
	if len(qb.UpdatedBindings) == 0 {
		errs = append(errs, errors.New("no updated fields provided for update operation"))
	}

	if len(qb.Having) != 0 {
		errs = append(errs, errors.New("HAVING is not allowed in UPDATE operations in the MySql driver"))
	}

	if qb.Offset > 0 {
		errs = append(errs, errors.New("OFFSET is not allowed in UPDATE in the MySql driver"))
	}

	if len(qb.GroupBy) > 0 {
		errs = append(errs, errors.New("GROUP BY is not allowed in UPDATE in the MySql driver"))
	}

	if len(qb.Unions) > 0 {
		errs = append(errs, errors.New("UNION is not allowed in UPDATE in the MySql driver"))
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
