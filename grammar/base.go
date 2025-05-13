package grammar

import (
	"github.com/iMohamedSheta/xqb/types"
)

// BaseGrammar provides common functionality for all grammars
type BaseGrammar struct{}

// GrammarInterface defines the methods that all grammars must implement
type GrammarInterface interface {
	CompileSelect(*types.QueryBuilderData) (string, []interface{}, error)
	compileBaseQuery(*types.QueryBuilderData) (string, []interface{}, error)
	compileSelectClause(*types.QueryBuilderData) (string, []interface{}, error)
	compileFromClause(*types.QueryBuilderData) (string, []interface{}, error)
	compileJoins(*types.QueryBuilderData) (string, []interface{}, error)
	compileWhereClause(*types.QueryBuilderData) (string, []interface{}, error)
	compileGroupByClause(*types.QueryBuilderData) (string, []interface{}, error)
	compileHavingClause(*types.QueryBuilderData) (string, []interface{}, error)
	compileOrderByClause(*types.QueryBuilderData) (string, []interface{}, error)
	compileLimitOffsetClause(*types.QueryBuilderData) (string, []interface{}, error)
	compileLockingClause(*types.QueryBuilderData) (string, []interface{}, error)
	compileCTEs(*types.QueryBuilderData) (string, []interface{}, error)

	CompileInsert(*types.QueryBuilderData) (string, []interface{}, error)
	CompileUpdate(*types.QueryBuilderData) (string, []interface{}, error)
	CompileDelete(*types.QueryBuilderData) (string, []interface{}, error)

	Build(qb *types.QueryBuilderData) (string, []interface{}, error)
}
