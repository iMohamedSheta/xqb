package grammar

// GetGrammar returns the appropriate grammar for the given driver
func GetGrammar(driverName string) GrammarInterface {
	switch driverName {
	case "mysql":
		return &MySQLGrammar{}
	default:
		return &MySQLGrammar{} // Default to MySQL grammar
	}
}
