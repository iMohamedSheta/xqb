package enums

// QueryType represents the type of query being built
type QueryType int

const (
	SELECT QueryType = 1
	INSERT QueryType = 2
	UPDATE QueryType = 3
	DELETE QueryType = 4

	CLAUSE_WHERE    QueryType = 5
	CLAUSE_SELECT   QueryType = 6
	CLAUSE_FROM     QueryType = 7
	CLAUSE_HAVING   QueryType = 8
	CLAUSE_ORDER_BY QueryType = 9
	CLAUSE_GROUP_BY QueryType = 10
	CLAUSE_LIMIT    QueryType = 11
	CLAUSE_OFFSET   QueryType = 12
	CLAUSE_UNION    QueryType = 13
	CLAUSE_JOIN     QueryType = 14
	CLAUSE_CTE      QueryType = 17
	CLAUSE_LOCKING  QueryType = 18
)
