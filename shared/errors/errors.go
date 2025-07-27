package errors

import "errors"

// Query Builder Errors
var (
	// ErrNotFound is returned when a query returns no results.
	// Commonly used to signal that a resource (row, record) doesn't exist.
	ErrNotFound = errors.New("xqb_resource_not_found")

	// ErrInvalidQuery is returned when query building fails due to syntax errors,
	// unsupported constructs, or invalid parameters.
	ErrInvalidQuery = errors.New("xqb_invalid_query")

	// ErrQueryFailed is returned when a query execution fails, often due to a database error.
	ErrQueryFailed = errors.New("xqb_query_failed")

	// ErrInvalidResult is returned when a query returns an unexpected result type or structure.
	// Commonly used to signal a mismatch between expected and actual result types.
	ErrInvalidResult = errors.New("xqb_invalid_result")

	// ErrUnsupportedFeature is returned when a feature is not supported by the underlying dialect,
	// such as streaming, chunking, or advanced Sql syntax.
	ErrUnsupportedFeature = errors.New("xqb_unsupported_feature")

	// ErrTransactionFailed is returned when a transaction could not be completed successfully,
	// often due to rollback or nested failure.
	ErrTransactionFailed = errors.New("xqb_transaction_failed")
)

// Database Manager Errors
var (
	// ErrNoConnection is returned when there is no database connection available.
	// Indicates a critical failure in establishing or maintaining a connection.
	ErrNoConnection = errors.New("xqb_no_database_connection")

	// ErrClosingConnection is returned when a database connection could not be closed.
	ErrClosingConnection = errors.New("xqb_failed_to_close_connection")
)
