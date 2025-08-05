package xqb

import (
	"database/sql"
	"errors"
)

type SqlQuery struct {
	connection string
	tx         *sql.Tx
	sql        string
	args       []any
	beforeExec func() error
	afterExec  func() error
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

func (s *SqlQuery) WithBeforeExec(f func() error) *SqlQuery {
	s.beforeExec = f
	return s
}

func (s *SqlQuery) WithAfterExec(f func() error) *SqlQuery {
	s.afterExec = f
	return s
}

// WithTx - set the transaction
func (s *SqlQuery) WithTx(tx *sql.Tx) *SqlQuery {
	s.tx = tx
	return s
}

// ExecuteSql - execute raw sql statement
func (s *SqlQuery) Execute() (sql.Result, error) {
	// Run before exec callback if set
	if s.beforeExec != nil {
		if err := s.beforeExec(); err != nil {
			return nil, err
		}
	}

	// If tx is set, use it to execute the sql statement
	if s.tx != nil {
		return s.tx.Exec(s.sql, s.args...)
	}

	// Otherwise, get the connection's sql.DB and execute the statement
	db, err := GetConnectionDB(s.connection)
	if err != nil {
		return nil, err
	}

	result, execErr := db.Exec(s.sql, s.args...)

	// Run after exec callback if set
	if s.afterExec != nil {
		if hookErr := s.afterExec(); hookErr != nil {
			if execErr == nil {
				return result, hookErr
			}
			return result, errors.Join(execErr, hookErr)
		}
	}

	return result, execErr
}

// QuerySql - query raw sql statement
func (s *SqlQuery) Query() (*sql.Rows, error) {
	// Run before exec callback if set
	if s.beforeExec != nil {
		if err := s.beforeExec(); err != nil {
			return nil, err
		}
	}

	// If tx is set, use it to query the sql statement
	if s.tx != nil {
		return s.tx.Query(s.sql, s.args...)
	}

	// Otherwise, get the connection's sql.DB and execute the statement
	db, err := GetConnectionDB(s.connection)
	if err != nil {
		return nil, err
	}

	rows, execErr := db.Query(s.sql, s.args...)

	// Run after exec callback if set
	if s.afterExec != nil {
		if hookErr := s.afterExec(); hookErr != nil {
			if execErr == nil {
				return rows, hookErr
			}
			return rows, errors.Join(execErr, hookErr)
		}
	}

	return rows, execErr
}

// QueryRow - query raw sql statement and scan the result into the pointed dest variable
// Example: Sql("SELECT * FROM users WHERE id = ?", 1).QueryRow(&user)
func (s *SqlQuery) QueryRow(dest ...any) error {
	if s.beforeExec != nil {
		if err := s.beforeExec(); err != nil {
			return err
		}
	}

	var row *sql.Row
	if s.tx != nil {
		row = s.tx.QueryRow(s.sql, s.args...)
	} else {
		db, err := GetConnectionDB(s.connection)
		if err != nil {
			return err
		}
		row = db.QueryRow(s.sql, s.args...)
	}

	scanErr := row.Scan(dest...)

	if s.afterExec != nil {
		if hookErr := s.afterExec(); hookErr != nil {
			if scanErr == nil {
				return hookErr
			}
			return errors.Join(scanErr, hookErr)
		}
	}

	return scanErr
}
