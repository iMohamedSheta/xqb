package grammar

import "github.com/iMohamedSheta/xqb/mysql"

// GetGrammar returns the appropriate grammar for the given driver
func GetGrammar(driverName string) GrammarInterface {
	switch driverName {
	case "mysql":
		return &mysql.MySQLGrammar{}
	default:
		return &mysql.MySQLGrammar{} // Default to MySQL grammar
	}
}
