package grammar

import "github.com/iMohamedSheta/xqb/shared/types"

// BaseGrammar provides common functionality for all grammars
type BaseGrammar struct{}

// GrammarInterface defines the methods that all grammars must implement
type GrammarInterface interface {
	CompileSelect(*types.QueryBuilderData) (string, []any, error)
	// compileBaseQuery(*types.QueryBuilderData) (string, []any, error)
	// compileSelectClause(*types.QueryBuilderData) (string, []any, error)
	// compileFromClause(*types.QueryBuilderData) (string, []any, error)
	// compileJoins(*types.QueryBuilderData) (string, []any, error)
	// compileWhereClause(*types.QueryBuilderData) (string, []any, error)
	// compileGroupByClause(*types.QueryBuilderData) (string, []any, error)
	// compileHavingClause(*types.QueryBuilderData) (string, []any, error)
	// compileOrderByClause(*types.QueryBuilderData) (string, []any, error)
	// compileLimitClause(*types.QueryBuilderData) (string, []any, error)
	// compileOffsetClause(*types.QueryBuilderData) (string, []any, error)
	// compileLockingClause(*types.QueryBuilderData) (string, []any, error)
	// compileCTEs(*types.QueryBuilderData) (string, []any, error)

	CompileInsert(*types.QueryBuilderData) (string, []any, error)
	CompileUpdate(*types.QueryBuilderData) (string, []any, error)
	CompileDelete(*types.QueryBuilderData) (string, []any, error)

	Build(qb *types.QueryBuilderData) (string, []any, error)
}
