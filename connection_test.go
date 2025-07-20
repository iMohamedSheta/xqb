package xqb_test

import (
	"database/sql"
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/stretchr/testify/assert"
)

func Test_ConnectionLifecycle(t *testing.T) {
	var db *sql.DB = nil
	name := "testdb"

	xqb.AddConnection(name, db)
	assert.True(t, xqb.HasConnection(name), "connection should exist after AddConnection")

	xqb.SetDefault(name)

	_, err := xqb.GetConnection(name)
	assert.Error(t, err, "getting nil connection should error")

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
	assert.NoError(t, err, "closing non-existent connection should not error")
}

func Test_CloseAllConnections(t *testing.T) {
	var db1 *sql.DB = nil
	var db2 *sql.DB = nil

	xqb.AddConnection("db1", db1)
	xqb.AddConnection("db2", db2)

	err := xqb.CloseAll()
	assert.NoError(t, err, "closing all connections should not error")
	assert.False(t, xqb.HasConnection("db1"), "nil connection db1 shouldn't remain after CloseAll")
	assert.False(t, xqb.HasConnection("db2"), "nil connection db2 shouldn't remain after CloseAll")
}
