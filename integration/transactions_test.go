package integration

import (
	"database/sql"
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTransaction_Commit db transaction success path.
func TestTransaction_Commit(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		err := xqb.TransactionOn(conn, func(tx *sql.Tx) error {
			if err := QB(conn, "xqb_users").WithTx(tx).Insert([]map[string]any{
				{"name": "TxOk", "email": "txok@example.com", "age": 40},
			}); err != nil {
				return err
			}
			_, err := QB(conn, "xqb_users").WithTx(tx).
				Where("name", "=", "Alice").
				Update(map[string]any{"age": 99})
			return err
		})
		require.NoError(t, err, "[%s] commit", conn)

		assert.Equal(t, int64(5), countUsers(t, conn))
		row, err := QB(conn, "xqb_users").Where("name", "=", "Alice").First()
		require.NoError(t, err)
		assert.Equal(t, int64(99), asInt64(row["age"]))
	})
}

// TestTransaction_Rollback transaction rollback on exception.
func TestTransaction_Rollback(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		err := xqb.TransactionOn(conn, func(tx *sql.Tx) error {
			if err := QB(conn, "xqb_users").WithTx(tx).Insert([]map[string]any{
				{"name": "TxFail", "email": "txfail@example.com"},
			}); err != nil {
				return err
			}
			// Force failure with an unknown column so the tx must roll back.
			_, err := QB(conn, "xqb_users").WithTx(tx).
				Where("name", "=", "Alice").
				Update(map[string]any{"no_such_column": 1})
			if err == nil {
				t.Logf("[%s] driver did not error on bad column inside tx; forcing rollback", conn)
				return sql.ErrTxDone
			}
			return err
		})
		require.Error(t, err, "[%s] expected rollback error", conn)

		// Neither the insert nor any partial update may survive.
		assert.Equal(t, int64(4), countUsers(t, conn))
		exists, err := QB(conn, "xqb_users").Where("email", "=", "txfail@example.com").Exists()
		require.NoError(t, err)
		assert.False(t, exists)
	})
}

// TestTransaction_DeleteRollback ensures deletes roll back too.
func TestTransaction_DeleteRollback(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		err := xqb.TransactionOn(conn, func(tx *sql.Tx) error {
			if _, err := QB(conn, "xqb_users").WithTx(tx).Where("status", "=", "active").Delete(); err != nil {
				return err
			}
			return sql.ErrTxDone // force rollback
		})
		require.Error(t, err)
		assert.Equal(t, int64(4), countUsers(t, conn), "[%s] rollback must restore rows", conn)
	})
}
