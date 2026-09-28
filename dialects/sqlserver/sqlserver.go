package sqlserver

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/iMohamedSheta/xqb/shared/enums"
	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// SqlServerDialect implements SQL Server (T-SQL) specific SQL syntax.
type SqlServerDialect struct {
}

func (d *SqlServerDialect) Getdialect() types.Dialect {
	return types.DialectSQLServer
}

// CompileSelect generates a SELECT SQL statement for SQL Server
func (d *SqlServerDialect) CompileSelect(qb *types.QueryBuilderData) (string, []any, error) {
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

// compileBaseQuery compiles a query without unions.
// NOTE: clause order differs from other dialects — OFFSET comes before
// FETCH (limit), matching SqlServerGrammar::$selectComponents.
func (d *SqlServerDialect) compileBaseQuery(qb *types.QueryBuilderData) (string, []any, error) {
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
		d.compileOffsetClause,
		d.compileLimitClause,
		d.compileLockingClause,
	}

	for _, compiler := range clauses {
		if err := d.AppendClause(&sql, &bindings, compiler, qb); err != nil {
			return "", nil, err
		}
	}

	return sql.String(), bindings, nil
}

func (d *SqlServerDialect) Build(qbd *types.QueryBuilderData) (string, []any, error) {
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

	return d.replaceQuestionMarksWithAtP(sql), bindings, nil
}

// AppendClause compiles and appends a clause to the SQL string and bindings
func (d *SqlServerDialect) AppendClause(sql *strings.Builder, bindings *[]any, compiler func(*types.QueryBuilderData) (string, []any, error), qb *types.QueryBuilderData) error {
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
func (d *SqlServerDialect) AppendError(qb *types.QueryBuilderData, err error) (string, []any, error) {
	qb.Errors = append(qb.Errors, err)
	return "", nil, err
}

// replaceQuestionMarksWithAtP converts ? placeholders to @p1, @p2, ...
// used by the Go SQL Server drivers (go-mssqldb).
func (d *SqlServerDialect) replaceQuestionMarksWithAtP(sql string) string {
	// Normalize any pre-existing $n / @pN placeholders back to ? first
	// (some raw fragments may already use them).
	reDollar := regexp.MustCompile(`\$\d+`)
	sql = reDollar.ReplaceAllString(sql, "?")
	reAtP := regexp.MustCompile(`(?i)@p\d+`)
	sql = reAtP.ReplaceAllString(sql, "?")

	// Then convert ? to @pN
	parts := strings.Split(sql, "?")
	if len(parts) == 1 {
		return sql
	}

	var b strings.Builder
	for i := 0; i < len(parts)-1; i++ {
		b.WriteString(parts[i])
		b.WriteString(fmt.Sprintf("@p%d", i+1))
	}
	b.WriteString(parts[len(parts)-1])

	return b.String()
}

func (d *SqlServerDialect) Wrap(value string) string {
	return wrapBracket(value)
}
