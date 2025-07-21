package dialects

import "github.com/iMohamedSheta/xqb/shared/types"

// BaseGrammar provides common functionality for all grammars
type BaseGrammar struct{}

// DialectInterface defines the methods that all grammars must implement
type DialectInterface interface {
	CompileSelect(*types.QueryBuilderData) (string, []any, error)
	CompileInsert(*types.QueryBuilderData) (string, []any, error)
	CompileUpdate(*types.QueryBuilderData) (string, []any, error)
	CompileDelete(*types.QueryBuilderData) (string, []any, error)

	Build(qb *types.QueryBuilderData) (string, []any, error)
}
