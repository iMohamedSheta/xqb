package dialects

import (
	"github.com/iMohamedSheta/xqb/dialects/mysql"
	"github.com/iMohamedSheta/xqb/dialects/postgres"
)

type Driver string

const (
	DriverMySQL    Driver = "mysql"
	DriverPostgres Driver = "postgres"
)

// GetDialect returns the appropriate grammar for the given driver
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
