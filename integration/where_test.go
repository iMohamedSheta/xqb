package integration

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Covers where family on every dialect.

func TestWhere_InNotIn(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		rows, err := QB(conn, "xqb_users").
			Select("name").WhereIn("name", []string{"Alice", "Carol"}).
			OrderBy("name", "ASC").Get()
		require.NoError(t, err)
		require.Len(t, rows, 2)

		rows, err = QB(conn, "xqb_users").
			Select("name").WhereNotIn("name", []string{"Alice", "Bob", "Carol"}).
			Get()
		require.NoError(t, err)
		require.Len(t, rows, 1)
		assert.Equal(t, "Dave", asString(rows[0]["name"]))
	})
}

func TestWhere_NullNotNull(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		rows, err := QB(conn, "xqb_users").Select("name").WhereNull("email").Get()
		require.NoError(t, err)
		require.Len(t, rows, 1)
		assert.Equal(t, "Dave", asString(rows[0]["name"]))

		rows, err = QB(conn, "xqb_users").Select("name").WhereNotNull("email").Get()
		require.NoError(t, err)
		require.Len(t, rows, 3)
	})
}

func TestWhere_Between(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		rows, err := QB(conn, "xqb_users").
			Select("name").WhereBetween("age", 26, 34).
			OrderBy("name", "ASC").Get()
		require.NoError(t, err)
		require.Len(t, rows, 1)
		assert.Equal(t, "Alice", asString(rows[0]["name"]))

		rows, err = QB(conn, "xqb_users").
			Select("name").WhereNotBetween("age", 26, 34).
			OrderBy("name", "ASC").Get()
		require.NoError(t, err)
		// Bob (25) + Carol (35); NULL row excluded by NOT BETWEEN semantics.
		require.Len(t, rows, 2)
	})
}

func TestWhere_OrGroups(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		rows, err := QB(conn, "xqb_users").
			Select("name").
			Where("status", "=", "inactive").
			OrWhere("name", "=", "Alice").
			OrderBy("name", "ASC").Get()
		require.NoError(t, err)
		require.Len(t, rows, 2)
		assert.Equal(t, "Alice", asString(rows[0]["name"]))
		assert.Equal(t, "Carol", asString(rows[1]["name"]))

		// Grouped: active AND (age > 26 OR name = Bob)
		rows, err = QB(conn, "xqb_users").
			Select("name").
			Where("status", "=", "active").
			WhereGroup(func(qb *xqb.QueryBuilder) {
				qb.Where("age", ">", 26).OrWhere("name", "=", "Bob")
			}).
			OrderBy("name", "ASC").Get()
		require.NoError(t, err)
		require.Len(t, rows, 2)
	})
}

func TestWhere_RawAndExists(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		rows, err := QB(conn, "xqb_users").
			Select("name").WhereRaw("age > ?", 30).
			OrderBy("name", "ASC").Get()
		require.NoError(t, err)
		require.Len(t, rows, 1)
		assert.Equal(t, "Carol", asString(rows[0]["name"]))

		// EXISTS: users that have at least one post.
		sub := QB(conn, "xqb_posts").Select("id").Where("xqb_posts.user_id", "=", 1)
		_ = sub
		rows, err = QB(conn, "xqb_users").
			Select("name").
			WhereExists(QB(conn, "xqb_posts").Select("id").WhereRaw("xqb_posts.user_id = xqb_users.id")).
			OrderBy("name", "ASC").Get()
		require.NoError(t, err)
		require.Len(t, rows, 2)
		assert.Equal(t, "Alice", asString(rows[0]["name"]))
		assert.Equal(t, "Bob", asString(rows[1]["name"]))

		// NOT EXISTS: users without posts.
		rows, err = QB(conn, "xqb_users").
			Select("name").
			WhereNotExists(QB(conn, "xqb_posts").Select("id").WhereRaw("xqb_posts.user_id = xqb_users.id")).
			OrderBy("name", "ASC").Get()
		require.NoError(t, err)
		require.Len(t, rows, 2)
	})
}

func TestWhere_InSubquery(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		rows, err := QB(conn, "xqb_users").
			Select("name").
			WhereInQuery("id", QB(conn, "xqb_posts").Select("user_id")).
			OrderBy("name", "ASC").Get()
		require.NoError(t, err)
		require.Len(t, rows, 2)
		assert.Equal(t, "Alice", asString(rows[0]["name"]))
		assert.Equal(t, "Bob", asString(rows[1]["name"]))
	})
}
