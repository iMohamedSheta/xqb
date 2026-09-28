package types

import "strings"

type Dialect string

const (
	DialectMySql     Dialect = "mysql"
	DialectMariaDB   Dialect = "mariadb"
	DialectPostgres  Dialect = "postgres"
	DialectSQLite    Dialect = "sqlite"
	DialectSQLServer Dialect = "sqlserver"
)

func (d Dialect) String() string {
	return string(d)
}

// Normalize maps common aliases (Go driver names)
// to the canonical dialect. Unknown values are lower-cased as-is.
func (d Dialect) Normalize() Dialect {
	switch Dialect(strings.ToLower(strings.TrimSpace(string(d)))) {
	case DialectMySql:
		return DialectMySql
	case DialectMariaDB:
		return DialectMariaDB
	case DialectPostgres, "pgsql", "postgresql", "pg":
		return DialectPostgres
	case DialectSQLite, "sqlite3":
		return DialectSQLite
	case DialectSQLServer, "sqlsrv", "mssql", "tsql":
		return DialectSQLServer
	default:
		return Dialect(strings.ToLower(strings.TrimSpace(string(d))))
	}
}
