package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUpdate_Basic update().
func TestUpdate_Basic(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		affected, err := QB(conn, "xqb_users").
			Where("name", "=", "Alice").
			Update(map[string]any{"status": "pending", "age": 31})
		require.NoError(t, err)
		assert.Equal(t, int64(1), affected)

		row, err := QB(conn, "xqb_users").Where("name", "=", "Alice").First()
		require.NoError(t, err)
		assert.Equal(t, "pending", asString(row["status"]))
		assert.Equal(t, int64(31), asInt64(row["age"]))
	})
}

// TestUpdate_MultipleRows mass update.
func TestUpdate_MultipleRows(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		affected, err := QB(conn, "xqb_users").
			Where("status", "=", "active").
			Update(map[string]any{"status": "verified"})
		require.NoError(t, err)
		assert.Equal(t, int64(2), affected)

		n, err := QB(conn, "xqb_users").Where("status", "=", "verified").Count("*")
		require.NoError(t, err)
		assert.Equal(t, int64(2), n)
	})
}

// TestUpdate_NoMatch returns 0 affected rows without error.
func TestUpdate_NoMatch(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		affected, err := QB(conn, "xqb_users").
			Where("name", "=", "Nobody").
			Update(map[string]any{"status": "x"})
		require.NoError(t, err)
		assert.Equal(t, int64(0), affected)
	})
}

// TestUpdate_InvalidColumn must error on every dialect.
func TestUpdate_InvalidColumn(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		_, err := QB(conn, "xqb_users").
			Where("id", "=", 1).
			Update(map[string]any{"no_such_column": "boom"})
		assert.Error(t, err, "[%s] expected error for unknown column", conn)
	})
}
