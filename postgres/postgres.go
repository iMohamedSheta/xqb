package postgres

import (
	"fmt"

	"github.com/iMohamedSheta/xqb/types"
)

// PostgresGrammar implements PostgreSQL-specific SQL syntax
type PostgresGrammar struct {
}

// CompileSelect generates a SELECT SQL statement for PostgreSQL
func (pg *PostgresGrammar) CompileSelect(qb *types.QueryBuilderData) (string, []interface{}, error) {
	return "", nil, nil
}

// CompileInsert generates an INSERT SQL statement for PostgreSQL
func (pg *PostgresGrammar) CompileInsert(qb *types.QueryBuilderData) (string, []interface{}, error) {
	return "", nil, nil
}

// CompileUpdate generates an UPDATE SQL statement for PostgreSQL
func (pg *PostgresGrammar) CompileUpdate(qb *types.QueryBuilderData) (string, []interface{}, error) {
	return "", nil, nil
}

// CompileDelete generates a DELETE SQL statement for PostgreSQL
func (pg *PostgresGrammar) CompileDelete(qb *types.QueryBuilderData) (string, []interface{}, error) {
	return "", nil, nil
}

// Build compiles the full SQL query for PostgreSQL
func (pg *PostgresGrammar) Build(qbd *types.QueryBuilderData) (string, []interface{}, error) {
	var sql string
	var bindings []interface{}
	var err error

	switch qbd.QueryType {
	case types.SELECT:
		sql, bindings, err = pg.CompileSelect(qbd)
	case types.INSERT:
		sql, bindings, err = pg.CompileInsert(qbd)
	case types.UPDATE:
		sql, bindings, err = pg.CompileUpdate(qbd)
	case types.DELETE:
		sql, bindings, err = pg.CompileDelete(qbd)
	}

	if err != nil {
		return "", nil, err
	}

	if sql == "" {
		return "", nil, fmt.Errorf("xqb_error_unexpected_error: couldn't build the query sql is empty")
	}

	return sql, bindings, nil
}
