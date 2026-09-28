package mariadb

import (
	"errors"
	"fmt"
	"strings"

	"github.com/iMohamedSheta/xqb/shared/enums"
	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/iMohamedSheta/xqb/shared/wrap"
)

// MariaDBDialect implements MariaDB-specific SQL syntax.
// MariaDB is wire-compatible with MySQL for the query builder
// (` quoting, ? placeholders, LIMIT/OFFSET, FOR UPDATE / LOCK IN SHARE MODE)
// and additionally supports INTERSECT / EXCEPT set operations (10.3+),
type MariaDBDialect struct {
}

func (d *MariaDBDialect) Getdialect() types.Dialect {
	return types.DialectMariaDB
}

// CompileSelect generates a SELECT SQL statement for MariaDB
func (d *MariaDBDialect) CompileSelect(qb *types.QueryBuilderData) (string, []any, error) {
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
func (d *MariaDBDialect) compileBaseQuery(qb *types.QueryBuilderData) (string, []any, error) {
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

func (d *MariaDBDialect) Build(qbd *types.QueryBuilderData) (string, []any, error) {
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
func (d *MariaDBDialect) AppendClause(sql *strings.Builder, bindings *[]any, compiler func(*types.QueryBuilderData) (string, []any, error), qb *types.QueryBuilderData) error {
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
func (d *MariaDBDialect) AppendError(qb *types.QueryBuilderData, err error) (string, []any, error) {
	qb.Errors = append(qb.Errors, err)
	return "", nil, err
}

func (d *MariaDBDialect) Wrap(value string) string {
	return wrap.Wrap(value, '`')
}
