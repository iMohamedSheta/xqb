package mysql

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/iMohamedSheta/xqb/shared/enums"
	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// MySQLDialect implements MySQL-specific SQL syntax
type MySQLDialect struct {
}

func (mg *MySQLDialect) GetDriver() types.Driver {
	return types.DriverMySQL
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

func (mg *MySQLDialect) Wrap(value string) string {
	value = strings.TrimSpace(value)
	lower := strings.ToLower(value)

	// Handle aliases (e.g., COUNT(id) AS total)
	if idx := strings.LastIndex(lower, " as "); idx != -1 {
		left := strings.TrimSpace(value[:idx])
		right := strings.TrimSpace(value[idx+4:])
		return fmt.Sprintf("%s AS %s", mg.Wrap(left), wrapMysqlValue(right))
	}

	// Don't wrap SQL expressions like COUNT(id)
	if mg.isLikelyExpr(lower) || mg.isLiteral(lower) {
		return value
	}

	// Handle dot notation like table.*
	if strings.HasSuffix(value, ".*") {
		parts := strings.SplitN(value, ".", 2)
		return wrapMysqlValue(parts[0]) + ".*"
	}

	// Handle dot notation like table.column
	segments := strings.Split(value, ".")
	for i := range segments {
		segments[i] = wrapMysqlValue(segments[i])
	}
	return strings.Join(segments, ".")
}

func wrapMysqlValue(value string) string {
	if value == "*" {
		return "*"
	}
	return "`" + strings.Trim(value, "`") + "`"
}

func (mg *MySQLDialect) isLikelyExpr(s string) bool {
	return strings.ContainsAny(s, "()+*/-")
}

func (mg *MySQLDialect) isLiteral(s string) bool {
	if s == "null" || s == "true" || s == "false" {
		return true
	}
	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return true // numeric literal
	}
	if strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'") {
		return true // string literal
	}
	return false
}
