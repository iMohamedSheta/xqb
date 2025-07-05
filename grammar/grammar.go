package grammar

import (
	"github.com/iMohamedSheta/xqb/mysql"
	"github.com/iMohamedSheta/xqb/postgres"
)

type Driver string

const (
	DriverMySQL    Driver = "mysql"
	DriverPostgres Driver = "postgres"
)

// GetGrammar returns the appropriate grammar for the given driver
func GetGrammar(driver Driver) GrammarInterface {
	switch driver {
	case DriverMySQL:
		return &mysql.MySQLGrammar{}
	case DriverPostgres:
		return &postgres.PostgresGrammar{}
	default:
		return &mysql.MySQLGrammar{} // Default to MySQL grammar
	}
}
