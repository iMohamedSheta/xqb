package enums

// QueryType represents the type of query being built
type QueryType int

const (
	SELECT QueryType = 1
	INSERT QueryType = 2
	UPDATE QueryType = 3
	DELETE QueryType = 4
)
