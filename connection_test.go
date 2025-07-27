package xqb_test

import (
	"database/sql"
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/stretchr/testify/assert"
)

func Test_ConnectionLifecycle(t *testing.T) {
	var db *sql.DB = nil
	name := "testdb"

	xqb.AddConnection(&xqb.Connection{
		Name:    name,
		DB:      db,
		Dialect: types.DialectMySql,
	})
	assert.False(t, xqb.HasConnection(name), "HasConnection should return false for nil connection")

	xqb.SetDefaultConnection(name)

	_, err := xqb.GetConnectionDB(name)
	assert.Error(t, err, "getting nil DB connection should error")

	_, err = xqb.GetConnection("")
	assert.Error(t, err, "getting default nil connection should error")

	err = xqb.Close(name)
	assert.NoError(t, err, "closing nil connection should not error")
	assert.False(t, xqb.HasConnection(name), "nil connection shouldn't exist after close")
}

func Test_ConnectionErrors(t *testing.T) {
	_, err := xqb.GetConnection("nonexistent")
	assert.Error(t, err, "getting non-existent connection should error")

	err = xqb.Close("nonexistent")
	assert.ErrorIs(t, err, errors.ErrNoConnection)
}

func Test_CloseAllConnections(t *testing.T) {
	old := xqb.DBManager()
	defer xqb.DBManager().SetConnections(old.GetConnections())

	var db1 *sql.DB = nil
	var db2 *sql.DB = nil

	xqb.AddConnection(&xqb.Connection{
		Name:    "db1",
		DB:      db1,
		Dialect: types.DialectMySql,
	})
	xqb.AddConnection(&xqb.Connection{
		Name:    "db2",
		DB:      db2,
		Dialect: types.DialectMySql,
	})

	err := xqb.CloseAll()
	assert.NoError(t, err, "closing all connections should not error")
	assert.False(t, xqb.HasConnection("db1"), "nil connection db1 shouldn't remain after CloseAll")
	assert.False(t, xqb.HasConnection("db2"), "nil connection db2 shouldn't remain after CloseAll")
}
