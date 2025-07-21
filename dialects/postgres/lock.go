package postgres

import (
	"strings"

	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileLockingClause compiles the locking clause
func (pg *PostgresDialect) compileLockingClause(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql strings.Builder

	// Check lock mode (FOR UPDATE, FOR SHARE, etc.)
	if lockVal, ok := qb.GetOption(types.OptionLock); ok {
		switch lockVal {
		case types.LockForUpdate:
			sql.WriteString(" FOR UPDATE")
		case types.LockInShare:
			sql.WriteString(" FOR SHARE")
		case types.LockNoKeyUpdate:
			sql.WriteString(" FOR NO KEY UPDATE")
		case types.LockKeyShare:
			sql.WriteString(" FOR KEY SHARE")
		}
	}

	// lock wait behavior (NOWAIT, SKIP LOCKED)
	if waitVal, ok := qb.GetOption(types.OptionLockWait); ok {
		switch waitVal {
		case types.LockNoWait:
			sql.WriteString(" NOWAIT")
		case types.LockSkipLocked:
			sql.WriteString(" SKIP LOCKED")
		}
	}

	return sql.String(), bindings, nil
}
