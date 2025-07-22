package postgres

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/iMohamedSheta/xqb/shared/enums"
	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// PostgresDialect implements Postgres-specific SQL syntax
type PostgresDialect struct {
}

func (pg *PostgresDialect) GetDriver() types.Driver {
	return types.DriverPostgres
}

// CompileSelect generates a SELECT SQL statement for Postgres
func (pg *PostgresDialect) CompileSelect(qb *types.QueryBuilderData) (string, []any, error) {
	if len(qb.Unions) == 0 {
		return pg.compileBaseQuery(qb)
	}

	var bindings []any
	var sql strings.Builder

	// Compile base SELECT
	baseSQL, baseBindings, err := pg.compileBaseQuery(qb)
	if err != nil {
		return "", nil, err
	}

	// Compile UNIONs
	unionSQL, unionBindings, err := pg.compileUnionClause(qb)
	if err != nil {
		return "", nil, err
	}

	// Wrap base query when unions exist
	sql.WriteString("(")
	sql.WriteString(baseSQL)
	sql.WriteString(")")

	// Append union part
	sql.WriteString(unionSQL)

	// Merge bindings
	if baseBindings != nil {
		bindings = append(bindings, baseBindings...)
	}
	if unionBindings != nil {
		bindings = append(bindings, unionBindings...)
	}

	// Check if there are any errors in building the query
	if len(qb.Errors) > 0 {
		errs := errors.Join(qb.Errors...)
		return "", nil, fmt.Errorf("%w: %s", xqbErr.ErrInvalidQuery, errs)
	}

	return sql.String(), bindings, nil
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

	return pg.replaceQuestionMarksWithDollar(sql), bindings, nil
}

// appendClause compiles and appends a clause to the SQL string and bindings
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
	// First we replace all $n with ? in the SQL string some sql is build with $n
	re := regexp.MustCompile(`\$\d+`)
	sql = re.ReplaceAllString(sql, "?")

	// Then we replace all ? with $n in the SQL string
	parts := strings.Split(sql, "?")
	if len(parts) == 1 {
		return sql
	}

	var b strings.Builder
	for i := 0; i < len(parts)-1; i++ {
		// Add the part
		b.WriteString(parts[i])
		// Add the $n in the end of the SQL string
		b.WriteString(fmt.Sprintf("$%d", i+1))
	}
	// Add the last part
	b.WriteString(parts[len(parts)-1])
	return b.String()
}
