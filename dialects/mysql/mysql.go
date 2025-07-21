package mysql

import (
	"errors"
	"fmt"
	"strings"

	"github.com/iMohamedSheta/xqb/shared/enums"
	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

const (
	Dialect = "mysql"
)

// MySQLDialect implements MySQL-specific SQL syntax
type MySQLDialect struct {
}

// CompileSelect generates a SELECT SQL statement for MySQL
func (mg *MySQLDialect) CompileSelect(qb *types.QueryBuilderData) (string, []any, error) {
	if len(qb.Unions) == 0 {
		return mg.compileBaseQuery(qb)
	}

	var bindings []any
	var sql strings.Builder

	// Compile base SELECT
	baseSQL, baseBindings, err := mg.compileBaseQuery(qb)
	if err != nil {
		return "", nil, err
	}

	// Compile UNIONs
	unionSQL, unionBindings, err := mg.compileUnionClause(qb)
	if err != nil {
		return "", nil, err
	}

	// Wrap base query when unions exist
	sql.WriteString("(")
	sql.WriteString(baseSQL)
	sql.WriteString(")")

	// Append union part
	sql.WriteString(unionSQL)

	// Combine bindings
	bindings = append(bindings, baseBindings...)
	bindings = append(bindings, unionBindings...)

	// Check if there are any errors in building the query
	if len(qb.Errors) > 0 {
		errs := errors.Join(qb.Errors...)
		return "", nil, fmt.Errorf("%w: %s", xqbErr.ErrInvalidQuery, errs)
	}

	return sql.String(), bindings, nil
}

// compileBaseQuery compiles a query without unions
func (mg *MySQLDialect) compileBaseQuery(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql strings.Builder

	// Compile each part of the query in order
	clauses := []func(*types.QueryBuilderData) (string, []any, error){
		mg.compileCTEs,
		mg.compileSelectClause,
		mg.compileFromClause,
		mg.compileJoins,
		mg.compileWhereClause,
		mg.compileGroupByClause,
		mg.compileHavingClause,
		mg.compileOrderByClause,
		mg.compileLimitClause,
		mg.compileOffsetClause,
		mg.compileLockingClause,
	}

	for _, compiler := range clauses {
		if err := appendClause(&sql, &bindings, compiler, qb); err != nil {
			return "", nil, err
		}
	}

	return sql.String(), bindings, nil
}

func (mg *MySQLDialect) Build(qbd *types.QueryBuilderData) (string, []any, error) {
	var sql string
	var bindings []any
	var err error

	switch qbd.QueryType {
	case enums.SELECT:
		sql, bindings, err = mg.CompileSelect(qbd)
	case enums.INSERT:
		sql, bindings, err = mg.CompileInsert(qbd)
	case enums.UPDATE:
		sql, bindings, err = mg.CompileUpdate(qbd)
	case enums.DELETE:
		sql, bindings, err = mg.CompileDelete(qbd)
	}

	if err != nil {
		return "", nil, err
	}

	if sql == "" {
		return "", nil, fmt.Errorf("%w: couldn't build the query sql is empty", xqbErr.ErrInvalidQuery)
	}

	return sql, bindings, nil
}

// appendClause compiles and appends a clause to the SQL string and bindings
func appendClause(sql *strings.Builder, bindings *[]any, compiler func(*types.QueryBuilderData) (string, []any, error), qb *types.QueryBuilderData) error {
	// compile clause closure
	part, partBindings, err := compiler(qb)
	if err != nil {
		return err
	}
	if part != "" {
		sql.WriteString(part)
	}
	if partBindings != nil {
		*bindings = append(*bindings, partBindings...)
	}
	return nil
}

// appendError appends an error to the query builder and returns it
func (mg *MySQLDialect) appendError(qb *types.QueryBuilderData, err error) (string, []any, error) {
	qb.Errors = append(qb.Errors, err)
	return "", nil, err
}
