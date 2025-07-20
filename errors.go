package xqb

import "errors"

var (
	// No results found
	ErrNotFound = errors.New("xqb - the resource was not found")

	// The table name was not set
	ErrMissingTable = errors.New("xqb - no table specified for query")

	// No columns selected
	ErrNoColumns = errors.New("xqb - no columns selected")

	// Placeholder count does not match bindings
	ErrPlaceholderMismatch = errors.New("xqb - placeholders and bindings count mismatch")

	// No database connection available
	ErrNoConnection = errors.New("xqb - database connection not available to execute query")

	// Query building failed due to invalid query
	ErrInvalidQuery = errors.New("xqb - invalid query build error")

	// Attempted to use chunking or streaming on an unsupported connection
	ErrUnsupportedFeature = errors.New("xqb - feature not supported by the driver")

	// The query was not valid and it executed
	ErrInvalidExecutedQuerySyntax = errors.New("xqb - invalid executed query syntax")

	// Unexpected row count
	ErrUnexpectedRowCount = errors.New("xqb - unexpected row count")

	// Invalid result type
	ErrInvalidResultType = errors.New("xqb - invalid result type")

	// Closing connection failed
	ErrClosingConnection = errors.New("xqb - closing connection failed")

	// Transaction failed
	ErrTransactionFailed = errors.New("xqb - transaction failed")
)
