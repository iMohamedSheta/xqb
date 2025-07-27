package xqb

import (
	"database/sql"
)

type SqlQuery struct {
	connection string
	tx         *sql.Tx
	sql        string
	args       []any
}

func Sql(sql string, args ...any) *SqlQuery {
	return &SqlQuery{
		connection: "default",
		sql:        sql,
		args:       args,
	}
}

// Connection - set the connection
func (s *SqlQuery) Connection(connection string) *SqlQuery {
	if connection == "" || !DBManager().HasConnection(connection) {
		connection = "default"
	}
	s.connection = connection
	return s
}

// WithTx - set the transaction
func (s *SqlQuery) WithTx(tx *sql.Tx) *SqlQuery {
	s.tx = tx
	return s
}

// ExecuteSql - execute raw sql statement
func (s *SqlQuery) Execute() (sql.Result, error) {
	if s.tx != nil {
		return s.tx.Exec(s.sql, s.args...)
	}

	db, err := GetConnectionDB(s.connection)
	if err != nil {
		return nil, err
	}

	return db.Exec(s.sql, s.args...)
}

// QuerySql - query raw sql statement
func (s *SqlQuery) Query() (*sql.Rows, error) {
	if s.tx != nil {
		return s.tx.Query(s.sql, s.args...)
	}

	db, err := GetConnectionDB(s.connection)
	if err != nil {
		return nil, err
	}

	return db.Query(s.sql, s.args...)
}

// QueryRowSql  - query raw sql statement and return one row
func (s *SqlQuery) QueryRow() (*sql.Row, error) {
	if s.tx != nil {
		return s.tx.QueryRow(s.sql, s.args...), nil
	}

	db, err := GetConnectionDB(s.connection)
	if err != nil {
		return nil, err
	}

	return db.QueryRow(s.sql, s.args...), nil
}
