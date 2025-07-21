package postgres

import (
	"fmt"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileFromClause compiles the FROM clause
func (mg *PostgresDialect) compileFromClause(qb *types.QueryBuilderData) (string, []any, error) {
	sql, bindings, err := mg.resolveTable(qb, "select", true)
	if err != nil || sql == "" {
		return "", bindings, err
	}

	return " FROM " + sql, bindings, nil
}

// resolveTable validates and returns the table or raw SQL used
func (mg *PostgresDialect) resolveTable(qb *types.QueryBuilderData, statement string, allowBindings bool) (string, []any, error) {
	if qb.Table == nil || (qb.Table.Raw == nil && qb.Table.Name == "") {
		if len(qb.WithCTEs) > 0 {
			return "", nil, nil
		}
		return mg.appendError(qb, fmt.Errorf("%w: table name is required for %s statement", xqbErr.ErrInvalidQuery, statement))
	}

	if qb.Table.Raw != nil && qb.Table.Name != "" {
		return mg.appendError(qb, fmt.Errorf("%w: both raw SQL and table name are set; choose one for %s statement", xqbErr.ErrInvalidQuery, statement))
	}

	if qb.Table.Raw != nil {
		if len(qb.Table.Raw.Bindings) > 0 && !allowBindings {
			return mg.appendError(qb, fmt.Errorf("%w: raw table cannot contain bindings in %s statement", xqbErr.ErrInvalidQuery, statement))
		}
		return qb.Table.Raw.SQL, qb.Table.Raw.Bindings, nil
	}

	return qb.Table.Name, nil, nil
}
