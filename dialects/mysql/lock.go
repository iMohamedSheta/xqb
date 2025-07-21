package mysql

import (
	"strings"

	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileLockingClause compiles the locking clause
func (mg *MySQLDialect) compileLockingClause(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql strings.Builder

	// Check lock mode
	if lockVal, ok := qb.GetOption(types.OptionLock); ok {
		switch lockVal {
		case types.LockForUpdate:
			sql.WriteString(" FOR UPDATE")
		case types.LockInShare:
			sql.WriteString(" LOCK IN SHARE MODE")
		}
	}

	// lock wait behavior (MySQL 8.0+)
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
