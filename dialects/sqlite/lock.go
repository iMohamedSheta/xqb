package sqlite

import (
	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileLockingClause compiles the locking clause.
// SQLite has no SELECT locking clause.
// Locks are silently ignored so shared query code keeps working.
func (d *SqliteDialect) compileLockingClause(qb *types.QueryBuilderData) (string, []any, error) {
	return "", nil, nil
}
