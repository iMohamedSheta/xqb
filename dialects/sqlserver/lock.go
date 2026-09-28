package sqlserver

import (
	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileLockingClause compiles the locking clause.
// Locks are compiled as WITH (ROWLOCK,...) hints on the FROM clause for
// SQL Server, so there is no trailing lock clause here.
func (d *SqlServerDialect) compileLockingClause(qb *types.QueryBuilderData) (string, []any, error) {
	return "", nil, nil
}
