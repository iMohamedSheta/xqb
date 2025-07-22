package dialects

import (
	"github.com/iMohamedSheta/xqb/dialects/mysql"
	"github.com/iMohamedSheta/xqb/dialects/postgres"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// GetDialect returns the appropriate dialect for the given driver
func GetDialect(driver types.Driver) DialectInterface {
	switch driver {
	case types.DriverMySQL:
		return &mysql.MySQLDialect{}
	case types.DriverPostgres:
		return &postgres.PostgresDialect{}
	default:
		return &mysql.MySQLDialect{} // Default to MySQL grammar
	}
}

// DialectInterface defines the methods that all grammars must implement
type DialectInterface interface {
	GetDriver() types.Driver
	CompileSelect(*types.QueryBuilderData) (string, []any, error)
	CompileInsert(*types.QueryBuilderData) (string, []any, error)
	CompileUpdate(*types.QueryBuilderData) (string, []any, error)
	CompileDelete(*types.QueryBuilderData) (string, []any, error)

	Build(qb *types.QueryBuilderData) (string, []any, error)
}
