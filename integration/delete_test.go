package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// delete().
func TestDelete_Basic(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		affected, err := QB(conn, "xqb_users").Where("name", "=", "Carol").Delete()
		require.NoError(t, err)
		assert.Equal(t, int64(1), affected)
		assert.Equal(t, int64(3), countUsers(t, conn))

		_, err = QB(conn, "xqb_users").Where("name", "=", "Carol").First()
		assert.Error(t, err, "[%s] deleted row must be gone", conn)
	})
}

// mass delete.
func TestDelete_MultipleRows(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		affected, err := QB(conn, "xqb_users").Where("status", "=", "active").Delete()
		require.NoError(t, err)
		assert.Equal(t, int64(2), affected)
		assert.Equal(t, int64(2), countUsers(t, conn))
	})
}

// TestDelete_NoMatch returns 0 without error.
func TestDelete_NoMatch(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		affected, err := QB(conn, "xqb_users").Where("name", "=", "Nobody").Delete()
		require.NoError(t, err)
		assert.Equal(t, int64(0), affected)
	})
}

// TestDelete_InvalidColumn must error (or affect 0 rows on sqlite).
//
// Note: SQLite treats double-quoted unknown identifiers in WHERE as string
// literals (MySQL-compat fallback), so
// `DELETE ... WHERE "no_such_column" = ?` matches 0 rows instead of
// erroring. Every other dialect raises an error.
func TestDelete_InvalidColumn(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		affected, err := QB(conn, "xqb_users").Where("no_such_column", "=", 1).Delete()
		if conn == "sqlite" {
			assert.NoError(t, err)
			assert.Equal(t, int64(0), affected)
			return
		}
		assert.Error(t, err, "[%s] expected error for unknown column", conn)
	})
}
