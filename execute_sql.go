package xqb

import (
	"context"
	"database/sql"
)

type SqlQuery struct {
	connection string
	tx         *sql.Tx
	sql        string
	args       []any
	ctx        context.Context // new field
	// beforeExec func(context.Context)
	afterExec func(context.Context)
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

func (s *SqlQuery) WithContext(ctx context.Context) *SqlQuery {
	s.ctx = ctx
	return s
}

// func (s *SqlQuery) WithBeforeExec(f func(ctx context.Context)) *SqlQuery {
// 	s.beforeExec = f
// 	return s
// }

func (s *SqlQuery) WithAfterExec(f func(ctx context.Context)) *SqlQuery {
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
	// if s.beforeExec != nil {
	// 	safeCall(func() {
	// 		s.beforeExec(s.ctx)
	// 	})
	// }

	var (
		result sql.Result
		err    error
	)

	if s.ctx == nil {
		s.ctx = context.Background()
	}

	if s.tx != nil {
		result, err = s.tx.ExecContext(s.ctx, s.sql, s.args...)
	} else {
		db, errConn := GetConnectionDB(s.connection)
		if errConn != nil {
			return nil, errConn
		}
		result, err = db.ExecContext(s.ctx, s.sql, s.args...)
	}

	if s.afterExec != nil {
		safeCall(func() {
			s.afterExec(s.ctx)
		})
	}

	return result, err
}

// QuerySql - query raw sql statement
func (s *SqlQuery) Query() (*sql.Rows, error) {
	// if s.beforeExec != nil {
	// 	safeCall(func() {
	// 		s.beforeExec(s.ctx)
	// 	})
	// }

	if s.ctx == nil {
		s.ctx = context.Background()
	}

	var (
		rows *sql.Rows
		err  error
	)

	if s.tx != nil {
		rows, err = s.tx.QueryContext(s.ctx, s.sql, s.args...)
	} else {
		db, errConn := GetConnectionDB(s.connection)
		if errConn != nil {
			return nil, errConn
		}
		rows, err = db.QueryContext(s.ctx, s.sql, s.args...)
	}

	if s.afterExec != nil {
		safeCall(func() {
			s.afterExec(s.ctx)
		})
	}

	return rows, err
}

// QueryRow - query raw sql statement and scan the result into the pointed dest variable
// Example: Sql("SELECT * FROM users WHERE id = ?", 1).QueryRow(&user)
func (s *SqlQuery) QueryRow(dest ...any) error {
	// if s.beforeExec != nil {
	// 	safeCall(func() {
	// 		s.beforeExec(s.ctx)
	// 	})
	// }

	if s.ctx == nil {
		s.ctx = context.Background()
	}

	var row *sql.Row

	if s.tx != nil {
		row = s.tx.QueryRowContext(s.ctx, s.sql, s.args...)
	} else {
		db, err := GetConnectionDB(s.connection)
		if err != nil {
			return err
		}
		row = db.QueryRowContext(s.ctx, s.sql, s.args...)
	}

	// don't return error before running afterExec hook
	err := row.Scan(dest...)

	if s.afterExec != nil {
		safeCall(func() {
			s.afterExec(s.ctx)
		})
	}

	return err
}

// Use safeCall to avoid panic if the callback function panics
func safeCall(f func()) {
	defer func() {
		_ = recover() // absorb panic
	}()
	f()
}
