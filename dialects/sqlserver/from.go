package sqlserver

import (
	"fmt"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileFromClause compiles the FROM clause.
// Row-level locks are expressed as WITH (ROWLOCK,...) table hints on
// SQL Server, not as a trailing FOR UPDATE clause.
func (d *SqlServerDialect) compileFromClause(qb *types.QueryBuilderData) (string, []any, error) {
	sql, bindings, err := d.resolveTable(qb, "select", true)
	if err != nil || sql == "" {
		return "", bindings, err
	}

	from := " FROM " + sql

	if lockVal, ok := qb.GetOption(types.OptionLock); ok {
		switch lockVal {
		case types.LockForUpdate, types.LockNoKeyUpdate:
			from += " WITH (ROWLOCK,UPDLOCK,HOLDLOCK)"
		case types.LockInShare, types.LockKeyShare:
			from += " WITH (ROWLOCK,HOLDLOCK)"
		}
	}

	return from, bindings, nil
}

// resolveTable validates and returns the table or raw SQL used
func (d *SqlServerDialect) resolveTable(qb *types.QueryBuilderData, statement string, allowBindings bool) (string, []any, error) {
	if qb.Table == nil || (qb.Table.Raw == nil && qb.Table.Name == "") {
		if len(qb.WithCTEs) > 0 {
			return "", nil, nil
		}
		return d.AppendError(qb, fmt.Errorf("%w: table name is required for %s statement", xqbErr.ErrInvalidQuery, statement))
	}

	if qb.Table.Raw != nil && qb.Table.Name != "" {
		return d.AppendError(qb, fmt.Errorf("%w: both raw SQL and table name are set; choose one for %s statement", xqbErr.ErrInvalidQuery, statement))
	}

	if qb.Table.Raw != nil {
		if len(qb.Table.Raw.Bindings) > 0 && !allowBindings {
			return d.AppendError(qb, fmt.Errorf("%w: raw table cannot contain bindings in %s statement", xqbErr.ErrInvalidQuery, statement))
		}
		return qb.Table.Raw.Sql, qb.Table.Raw.Bindings, nil
	}

	return d.Wrap(qb.Table.Name), nil, nil
}
