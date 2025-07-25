package mysql

import (
	"errors"
	"fmt"
	"strings"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

func (mg *MySQLDialect) CompileDelete(qb *types.QueryBuilderData) (string, []any, error) {
	tableName, _, err := mg.resolveTable(qb, "delete", false)
	if err != nil {
		return "", nil, err
	}

	// validate query builder delete build
	if err := mg.validateDelete(qb); err != nil {
		return "", nil, err
	}

	var bindings []any
	var sql strings.Builder

	if len(qb.DeleteFrom) > 0 {
		sql.WriteString(fmt.Sprintf("DELETE %s FROM %s", strings.Join(qb.DeleteFrom, ", "), tableName))
	} else {
		sql.WriteString(fmt.Sprintf("DELETE FROM %s", tableName))
	}

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

// ValidateDelete validates the delete operation for postgres driver
func (mg *MySQLDialect) validateDelete(qb *types.QueryBuilderData) error {
	var errs []error

	if len(qb.Where) == 0 {
		errs = append(errs, errors.New("DELETE without WHERE clause is dangerous; we don't allow that"))
	}

	if len(qb.OrderBy) > 0 && qb.Limit == 0 {
		errs = append(errs, errors.New("ORDER BY in DELETE is only allowed with LIMIT in the MySQL driver"))
	}

	if len(qb.Having) != 0 {
		errs = append(errs, errors.New("HAVING is not allowed in DELETE in the MySQL driver"))
	}

	if qb.Offset > 0 {
		errs = append(errs, errors.New("OFFSET is not allowed in DELETE in the MySQL driver"))
	}

	if len(qb.GroupBy) > 0 {
		errs = append(errs, errors.New("GROUP BY is not allowed in DELETE in the MySQL driver"))
	}

	if len(qb.Unions) > 0 {
		errs = append(errs, errors.New("UNION is not allowed in DELETE in the MySQL driver"))
	}

	if len(qb.Columns) > 0 {
		errs = append(errs, errors.New("COLUMNS are not valid in DELETE queries in the MySQL driver"))
	}

	if qb.Distinct || qb.IsUsingDistinct {
		errs = append(errs, errors.New("DISTINCT is not valid in DELETE queries in the MySQL driver"))
	}

	if len(qb.Options) != 0 {
		for option := range qb.Options {
			errs = append(errs, errors.New(option.String()+" Options are Not supported in DELETE queries"))
		}
	}

	if len(qb.Joins) > 0 && len(qb.DeleteFrom) == 0 {
		errs = append(errs, errors.New("DELETE with JOINs requires specifying which table(s) to delete from (via Delete({tables...}))"))
	}

	if len(errs) != 0 {
		return fmt.Errorf("%w: %s", xqbErr.ErrInvalidQuery, errors.Join(errs...))
	}

	return nil
}
