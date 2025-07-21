package dialects

import (
	"github.com/iMohamedSheta/xqb/dialects/mysql"
	"github.com/iMohamedSheta/xqb/dialects/postgres"
	"github.com/iMohamedSheta/xqb/shared/types"
)

type Driver string

const (
	DriverMySQL    Driver = "mysql"
	DriverPostgres Driver = "postgres"
)

// GetDialect returns the appropriate dialect for the given driver
func GetDialect(driver Driver) DialectInterface {
	switch driver {
	case DriverMySQL:
		return &mysql.MySQLDialect{}
	case DriverPostgres:
		return &postgres.PostgresDialect{}
	default:
		return &mysql.MySQLDialect{} // Default to MySQL grammar
	}
}

// DialectInterface defines the methods that all grammars must implement
type DialectInterface interface {
	CompileSelect(*types.QueryBuilderData) (string, []any, error)
	CompileInsert(*types.QueryBuilderData) (string, []any, error)
	CompileUpdate(*types.QueryBuilderData) (string, []any, error)
	CompileDelete(*types.QueryBuilderData) (string, []any, error)

	Build(qb *types.QueryBuilderData) (string, []any, error)
}
