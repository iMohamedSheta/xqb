package dialects

import (
	"github.com/iMohamedSheta/xqb/dialects/mysql"
	"github.com/iMohamedSheta/xqb/dialects/postgres"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// GetDialect returns the appropriate dialect for the given dialect
func GetDialect(dialect types.Dialect) DialectInterface {
	switch dialect {
	case types.DialectMySql:
		return &mysql.MySqlDialect{}
	case types.DialectPostgres:
		return &postgres.PostgresDialect{}
	default:
		return &mysql.MySqlDialect{} // Default to MySql grammar
	}
}

// DialectInterface defines the methods that all grammars must implement
type DialectInterface interface {
	Getdialect() types.Dialect
	CompileSelect(*types.QueryBuilderData) (string, []any, error)
	CompileInsert(*types.QueryBuilderData) (string, []any, error)
	CompileUpdate(*types.QueryBuilderData) (string, []any, error)
	CompileDelete(*types.QueryBuilderData) (string, []any, error)

	Build(qb *types.QueryBuilderData) (string, []any, error)
}
