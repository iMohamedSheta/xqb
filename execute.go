package xqb

import "database/sql"

// Sql - execute raw sql statement
func ExecuteSql(sql string, args ...any) (sql.Result, error) {
	dbManager := GetDBManager()

	if !dbManager.IsDBConnected() {
		return nil, ErrNoConnection
	}

	db, err := dbManager.GetDB()
	if err != nil {
		return nil, err
	}

	return db.Exec(sql, args...)
}

// QuerySql - query raw sql statement
func QuerySql(sql string, args ...any) (*sql.Rows, error) {
	dbManager := GetDBManager()
	if !dbManager.IsDBConnected() {
		return nil, ErrNoConnection
	}
	db, err := dbManager.GetDB()
	if err != nil {
		return nil, err
	}
	return db.Query(sql, args...)
}

// QueryRowSql  - query raw sql statement and return one row
func QueryRowSql(sql string, args ...any) (*sql.Row, error) {
	dbManager := GetDBManager()
	if !dbManager.IsDBConnected() {
		return nil, ErrNoConnection
	}
	db, err := dbManager.GetDB()
	if err != nil {
		return nil, err
	}
	return db.QueryRow(sql, args...), nil
}
