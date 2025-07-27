package postgres

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/iMohamedSheta/xqb/shared/enums"
	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/iMohamedSheta/xqb/shared/wrap"
)

// PostgresDialect implements Postgres-specific Sql syntax
type PostgresDialect struct {
}

func (d *PostgresDialect) Getdialect() types.Dialect {
	return types.DialectPostgres
}

// CompileSelect generates a SELECT Sql statement for Postgres
func (d *PostgresDialect) CompileSelect(qb *types.QueryBuilderData) (string, []any, error) {
	if len(qb.Unions) == 0 {
		return d.compileBaseQuery(qb)
	}

	var bindings []any
	var sql string

	// Compile base SELECT
	baseSql, baseBindings, err := d.compileBaseQuery(qb)
	if err != nil {
		return "", nil, err
	}

	// Compile UNIONs
	unionSql, unionBindings, err := d.compileUnionClause(qb)
	if err != nil {
		return "", nil, err
	}

	// Wrap base query when unions exist then append union part
	sql += "(" + baseSql + ")" + unionSql

	// Merge bindings
	bindings = append(bindings, baseBindings...)
	bindings = append(bindings, unionBindings...)

	return sql, bindings, nil
}

// compileBaseQuery compiles a query without unions
func (d *PostgresDialect) compileBaseQuery(qb *types.QueryBuilderData) (string, []any, error) {
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

func (d *PostgresDialect) Build(qbd *types.QueryBuilderData) (string, []any, error) {
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

	return d.replaceQuestionMarksWithDollar(sql), bindings, nil
}

// appendClause compiles and appends a clause to the Sql string and bindings
func (d *PostgresDialect) AppendClause(sql *strings.Builder, bindings *[]any, compiler func(*types.QueryBuilderData) (string, []any, error), qb *types.QueryBuilderData) error {
	// compile clause closure
	part, partBindings, err := compiler(qb)
	if err != nil {
		return err
	}

	sql.WriteString(part)

	if partBindings != nil {
		*bindings = append(*bindings, partBindings...)
	}
	return nil
}

// appendError appends an error to the query builder and returns it
func (d *PostgresDialect) AppendError(qb *types.QueryBuilderData, err error) (string, []any, error) {
	qb.Errors = append(qb.Errors, err)
	return "", nil, err
}

func (d *PostgresDialect) replaceQuestionMarksWithDollar(sql string) string {
	// First we replace all $n with ? in the Sql string some sql is build with $n
	re := regexp.MustCompile(`\$\d+`)
	sql = re.ReplaceAllString(sql, "?")

	// Then we replace all ? with $n in the Sql string
	parts := strings.Split(sql, "?")
	if len(parts) == 1 {
		return sql
	}

	var b string
	for i := 0; i < len(parts)-1; i++ {
		// Add the part then Add the $n in the end of the Sql string
		b += parts[i] + fmt.Sprintf("$%d", i+1)
	}
	// Add the last part
	b += parts[len(parts)-1]

	return b
}
func (d *PostgresDialect) Wrap(value string) string {
	return wrap.Wrap(value, '"')
}
