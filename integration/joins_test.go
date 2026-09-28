package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestJoin_Inner join().
func TestJoin_Inner(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		rows, err := QB(conn, "xqb_users").
			Select("xqb_users.name", "xqb_posts.title").
			Join("xqb_posts", "xqb_users.id = xqb_posts.user_id").
			OrderBy("xqb_posts.id", "ASC").
			Get()
		require.NoError(t, err)
		require.Len(t, rows, 3)
		assert.Equal(t, "Alice", asString(rows[0]["name"]))
	})
}

// TestJoin_Left leftJoin(): users without posts appear with NULL.
func TestJoin_Left(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		rows, err := QB(conn, "xqb_users").
			Select("xqb_users.name", "xqb_posts.title").
			LeftJoin("xqb_posts", "xqb_users.id = xqb_posts.user_id").
			OrderBy("xqb_users.id", "ASC").
			OrderBy("xqb_posts.id", "ASC").
			Get()
		require.NoError(t, err)
		// Alice x2, Bob x1, Carol x1 (null post), Dave x1 (null post) = 5
		require.Len(t, rows, 5, "[%s] left join row count", conn)
	})
}

// TestJoin_WithWhere join + where combo.
func TestJoin_WithWhere(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		rows, err := QB(conn, "xqb_users").
			Select("xqb_users.name").
			Join("xqb_posts", "xqb_users.id = xqb_posts.user_id").
			Where("xqb_users.name", "=", "Alice").
			OrderBy("xqb_posts.id", "ASC").
			Get()
		require.NoError(t, err)
		require.Len(t, rows, 2)
	})
}

// TestUnion union(): active + inactive names.
//
// Note: xqb wraps compound selects as `(SELECT ...) UNION (SELECT ...)`.
// SQLite rejects parenthesized compound selects (syntax error near "("),
// so execution is skipped there — ToSql coverage lives in the unit tests.
func TestUnion(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		if conn == "sqlite" {
			t.Skip("parenthesized UNION not executable on sqlite")
		}
		resetTables(t, conn)

		q2 := QB(conn, "xqb_users").Select("name").Where("status", "=", "inactive")
		rows, err := QB(conn, "xqb_users").
			Select("name").Where("status", "=", "active").
			Union(q2).
			Get()
		require.NoError(t, err)
		require.Len(t, rows, 3)
	})
}

// TestUnionAll unionAll(): duplicates preserved.
func TestUnionAll(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		if conn == "sqlite" {
			t.Skip("parenthesized UNION ALL not executable on sqlite")
		}
		resetTables(t, conn)

		rows, err := QB(conn, "xqb_users").
			Select("status").WhereNotNull("status").
			UnionAll(QB(conn, "xqb_users").Select("status").Where("name", "=", "Alice")).
			Get()
		require.NoError(t, err)
		// 3 non-null statuses + 1 extra Alice row = 4
		require.Len(t, rows, 4)
	})
}
