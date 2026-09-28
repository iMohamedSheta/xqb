package sqlite

import (
	"errors"
	"fmt"
	"strings"

	"github.com/iMohamedSheta/xqb/shared/enums"
	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/iMohamedSheta/xqb/shared/wrap"
)

// SqliteDialect implements SQLite-specific SQL syntax.
// " quoting, ? placeholders, LIMIT/OFFSET, no locking clause,
// UNION/INTERSECT/EXCEPT support, ON CONFLICT upsert and RETURNING support.
type SqliteDialect struct {
}

func (d *SqliteDialect) Getdialect() types.Dialect {
	return types.DialectSQLite
}

// CompileSelect generates a SELECT SQL statement for SQLite
func (d *SqliteDialect) CompileSelect(qb *types.QueryBuilderData) (string, []any, error) {
	if len(qb.Unions) == 0 {
		return d.compileBaseQuery(qb)
	}

	var bindings []any
	var sql string

	// Compile base SELECT
	baseSQL, baseBindings, err := d.compileBaseQuery(qb)
	if err != nil {
		return "", nil, err
	}

	// Compile UNIONs
	unionSQL, unionBindings, err := d.compileUnionClause(qb)
	if err != nil {
		return "", nil, err
	}

	// Wrap base query when unions exist then append union part
	sql += "(" + baseSQL + ")" + unionSQL

	// Combine bindings
	bindings = append(bindings, baseBindings...)
	bindings = append(bindings, unionBindings...)

	return sql, bindings, nil
}

// compileBaseQuery compiles a query without unions
func (d *SqliteDialect) compileBaseQuery(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql strings.Builder

	// Compile each part of the query in order
	clauses := []func(*types.QueryBuilderData) (string, []any, error){
		d.compileCTEs,
		d.compileSelectClause,
		d.compileFromClause,
		d.compileJoins,
		d.compileWhereClause,
		d.compileGroupByClause,
		d.compileHavingClause,
		d.compileOrderByClause,
		d.compileLimitClause,
		d.compileOffsetClause,
		d.compileLockingClause,
	}

	for _, compiler := range clauses {
		if err := d.AppendClause(&sql, &bindings, compiler, qb); err != nil {
			return "", nil, err
		}
	}

	return sql.String(), bindings, nil
}

func (d *SqliteDialect) Build(qbd *types.QueryBuilderData) (string, []any, error) {
	var sql string
	var bindings []any
	var err error

	switch qbd.QueryType {
	case enums.SELECT:
		sql, bindings, err = d.CompileSelect(qbd)
	case enums.INSERT:
		sql, bindings, err = d.CompileInsert(qbd)
	case enums.UPDATE:
		sql, bindings, err = d.CompileUpdate(qbd)
	case enums.DELETE:
		sql, bindings, err = d.CompileDelete(qbd)
	}

	if err != nil {
		return "", nil, err
	}

	if sql == "" {
		return "", nil, fmt.Errorf("%w: couldn't build the query sql is empty", xqbErr.ErrInvalidQuery)
	}

	// Check if there are any errors in building the query
	if len(qbd.Errors) > 0 {
		return "", nil, errors.Join(qbd.Errors...)
	}

	return sql, bindings, nil
}

// AppendClause compiles and appends a clause to the SQL string and bindings
func (d *SqliteDialect) AppendClause(sql *strings.Builder, bindings *[]any, compiler func(*types.QueryBuilderData) (string, []any, error), qb *types.QueryBuilderData) error {
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

// AppendError appends an error to the query builder and returns it
func (d *SqliteDialect) AppendError(qb *types.QueryBuilderData, err error) (string, []any, error) {
	qb.Errors = append(qb.Errors, err)
	return "", nil, err
}

func (d *SqliteDialect) Wrap(value string) string {
	return wrap.Wrap(value, '"')
}
