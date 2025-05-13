package grammar

import (
	"strings"

	"github.com/iMohamedSheta/xqb/types"
)

// compileLockingClause compiles the locking clause
func (mg *MySQLGrammar) compileLockingClause(qb *types.QueryBuilderData) (string, []interface{}, error) {
	var bindings []interface{}
	var sql strings.Builder

	if qb.IsForUpdate {
		sql.WriteString(" FOR UPDATE")
	} else if qb.IsLockInShareMode {
		sql.WriteString(" LOCK IN SHARE MODE")
	}

	return sql.String(), bindings, nil
}
