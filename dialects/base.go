package dialects

import (
	"github.com/iMohamedSheta/xqb/dialects/mariadb"
	"github.com/iMohamedSheta/xqb/dialects/mysql"
	"github.com/iMohamedSheta/xqb/dialects/postgres"
	"github.com/iMohamedSheta/xqb/dialects/sqlite"
	"github.com/iMohamedSheta/xqb/dialects/sqlserver"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// GetDialect returns the appropriate dialect for the given dialect.
// Aliases (pgsql, sqlite3, sqlsrv, mssql, ...) are normalized
func GetDialect(dialect types.Dialect) DialectInterface {
	switch dialect.Normalize() {
	case types.DialectMySql:
		return &mysql.MySqlDialect{}
	case types.DialectMariaDB:
		return &mariadb.MariaDBDialect{}
	case types.DialectPostgres:
		return &postgres.PostgresDialect{}
	case types.DialectSQLite:
		return &sqlite.SqliteDialect{}
	case types.DialectSQLServer:
		return &sqlserver.SqlServerDialect{}
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
