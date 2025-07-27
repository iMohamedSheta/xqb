package postgres

import (
	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileLockingClause compiles the locking clause
func (d *PostgresDialect) compileLockingClause(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql string

	// Check lock mode (FOR UPDATE, FOR SHARE, etc.)
	if lockVal, ok := qb.GetOption(types.OptionLock); ok {
		switch lockVal {
		case types.LockForUpdate:
			sql += " FOR UPDATE"
		case types.LockInShare:
			sql += " FOR SHARE"
		case types.LockNoKeyUpdate:
			sql += " FOR NO KEY UPDATE"
		case types.LockKeyShare:
			sql += " FOR KEY SHARE"
		}
	}

	// lock wait behavior (NOWAIT, SKIP LOCKED)
	if waitVal, ok := qb.GetOption(types.OptionLockWait); ok {
		switch waitVal {
		case types.LockNoWait:
			sql += " NOWAIT"
		case types.LockSkipLocked:
			sql += " SKIP LOCKED"
		}
	}

	return sql, bindings, nil
}
