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

func (pg *PostgresDialect) GetDriver() types.Driver {
	return types.DriverPostgres
}

// CompileSelect generates a SELECT Sql statement for Postgres
func (pg *PostgresDialect) CompileSelect(qb *types.QueryBuilderData) (string, []any, error) {
	if len(qb.Unions) == 0 {
		return pg.compileBaseQuery(qb)
	}

	var bindings []any
	var sql string

	// Compile base SELECT
	baseSql, baseBindings, err := pg.compileBaseQuery(qb)
	if err != nil {
		return "", nil, err
	}

	// Compile UNIONs
	unionSql, unionBindings, err := pg.compileUnionClause(qb)
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
func (pg *PostgresDialect) compileBaseQuery(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql strings.Builder

	// Compile each part of the query in order
	clauses := []func(*types.QueryBuilderData) (string, []any, error){
		pg.compileCTEs,
		pg.compileSelectClause,
		pg.compileFromClause,
		pg.compileJoins,
		pg.compileWhereClause,
		pg.compileGroupByClause,
		pg.compileHavingClause,
		pg.compileOrderByClause,
		pg.compileLimitClause,
		pg.compileOffsetClause,
		pg.compileLockingClause,
	}

	for _, compiler := range clauses {
		if err := appendClause(&sql, &bindings, compiler, qb); err != nil {
			return "", nil, err
		}
	}

	return sql.String(), bindings, nil
}

func (pg *PostgresDialect) Build(qbd *types.QueryBuilderData) (string, []any, error) {
	var sql string
	var bindings []any
	var err error

	switch qbd.QueryType {
	case enums.SELECT:
		sql, bindings, err = pg.CompileSelect(qbd)
	case enums.INSERT:
		sql, bindings, err = pg.CompileInsert(qbd)
	case enums.UPDATE:
		sql, bindings, err = pg.CompileUpdate(qbd)
	case enums.DELETE:
		sql, bindings, err = pg.CompileDelete(qbd)
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

	return pg.replaceQuestionMarksWithDollar(sql), bindings, nil
}

// appendClause compiles and appends a clause to the Sql string and bindings
func appendClause(sql *strings.Builder, bindings *[]any, compiler func(*types.QueryBuilderData) (string, []any, error), qb *types.QueryBuilderData) error {
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
func (pg *PostgresDialect) appendError(qb *types.QueryBuilderData, err error) (string, []any, error) {
	qb.Errors = append(qb.Errors, err)
	return "", nil, err
}

func (pg *PostgresDialect) replaceQuestionMarksWithDollar(sql string) string {
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
func (pg *PostgresDialect) Wrap(value string) string {
	return wrap.Wrap(value, '"')
}
