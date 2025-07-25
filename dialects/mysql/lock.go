package mysql

import (
	"fmt"
	"strings"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
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
		default:
			return "", nil, fmt.Errorf("%w: invalid lock mode %q for MySQL dialect", xqbErr.ErrInvalidQuery, lockVal)
		}

		// lock wait behavior (MySQL 8.0+)
		if waitVal, ok := qb.GetOption(types.OptionLockWait); ok {
			switch waitVal {
			case types.LockNoWait:
				sql.WriteString(" NOWAIT")
			case types.LockSkipLocked:
				sql.WriteString(" SKIP LOCKED")
			default:
				return "", nil, fmt.Errorf("%w: invalid lock wait behavior %q for MySQL dialect", xqbErr.ErrInvalidQuery, waitVal)
			}
		}
	}

	return sql.String(), bindings, nil
}
