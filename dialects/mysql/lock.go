package mysql

import (
	"strings"

	"github.com/iMohamedSheta/xqb/shared/types"
)

// compileLockingClause compiles the locking clause
func (mg *MySQLGrammar) compileLockingClause(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql strings.Builder

	if qb.IsLockedForUpdate {
		sql.WriteString(" FOR UPDATE")
	} else if qb.IsInSharedLock {
		sql.WriteString(" LOCK IN SHARE MODE")
	}

	return sql.String(), bindings, nil
}
