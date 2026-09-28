package mariadb

import (
	"fmt"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileLockingClause compiles the locking clause
func (d *MariaDBDialect) compileLockingClause(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql string

	// Check lock mode
	if lockVal, ok := qb.GetOption(types.OptionLock); ok {
		switch lockVal {
		case types.LockForUpdate:
			sql += " FOR UPDATE"
		case types.LockInShare:
			sql += " LOCK IN SHARE MODE"
		default:
			return "", nil, fmt.Errorf("%w: invalid lock mode %q for MariaDB dialect", xqbErr.ErrInvalidQuery, lockVal)
		}

		// lock wait behavior (MariaDB 10.6+)
		if waitVal, ok := qb.GetOption(types.OptionLockWait); ok {
			switch waitVal {
			case types.LockNoWait:
				sql += " NOWAIT"
			case types.LockSkipLocked:
				sql += " SKIP LOCKED"
			default:
				return "", nil, fmt.Errorf("%w: invalid lock wait behavior %q for MariaDB dialect", xqbErr.ErrInvalidQuery, waitVal)
			}
		}
	}

	return sql, bindings, nil
}
