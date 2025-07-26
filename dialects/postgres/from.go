package postgres

import (
	"fmt"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileFromClause compiles the FROM clause
func (pg *PostgresDialect) compileFromClause(qb *types.QueryBuilderData) (string, []any, error) {
	sql, bindings, err := pg.resolveTable(qb, "select", true)
	if err != nil || sql == "" {
		return "", bindings, err
	}

	return " FROM " + sql, bindings, nil
}

// resolveTable validates and returns the table or raw Sql used
func (pg *PostgresDialect) resolveTable(qb *types.QueryBuilderData, statement string, allowBindings bool) (string, []any, error) {
	if qb.Table == nil || (qb.Table.Raw == nil && qb.Table.Name == "") {
		if len(qb.WithCTEs) > 0 {
			return "", nil, nil
		}
		return pg.appendError(qb, fmt.Errorf("%w: table name is required for %s statement", xqbErr.ErrInvalidQuery, statement))
	}

	if qb.Table.Raw != nil && qb.Table.Name != "" {
		return pg.appendError(qb, fmt.Errorf("%w: both raw Sql and table name are set; choose one for %s statement", xqbErr.ErrInvalidQuery, statement))
	}

	if qb.Table.Raw != nil {
		if len(qb.Table.Raw.Bindings) > 0 && !allowBindings {
			return pg.appendError(qb, fmt.Errorf("%w: raw table cannot contain bindings in %s statement", xqbErr.ErrInvalidQuery, statement))
		}
		return qb.Table.Raw.Sql, qb.Table.Raw.Bindings, nil
	}

	return pg.Wrap(qb.Table.Name), nil, nil
}
