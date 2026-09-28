package integration

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGroupByHaving mirrors  groupBy()/having().
func TestGroupByHaving(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		rows, err := QB(conn, "xqb_users").
			Select("status", xqb.Count("*", "total")).
			WhereNotNull("status").
			GroupBy("status").
			OrderBy("status", "ASC").
			Get()
		require.NoError(t, err)
		require.Len(t, rows, 2)

		// HAVING COUNT(*) > 1 keeps only 'active' (2 rows).
		// HavingRaw is used (instead of alias reference) for max
		// cross-dialect compatibility.
		rows, err = QB(conn, "xqb_users").
			Select("status", xqb.Count("*", "total")).
			WhereNotNull("status").
			GroupBy("status").
			HavingRaw("COUNT(*) > 1").
			Get()
		require.NoError(t, err)
		require.Len(t, rows, 1)
		assert.Equal(t, "active", asString(rows[0]["status"]))
	})
}

// CTE over active users.
func TestCTE(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		cte := QB(conn, "xqb_users").Select("id", "name").Where("status", "=", "active")
		rows, err := QB(conn, "active_users").
			With("active_users", cte).
			Select("name").
			OrderBy("name", "ASC").
			Get()
		require.NoError(t, err)
		require.Len(t, rows, 2)
		assert.Equal(t, "Alice", asString(rows[0]["name"]))
	})
}

func TestSubqueryFrom(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		sub := QB(conn, "xqb_users").Select("id", "name").Where("status", "=", "active")
		rows, err := QB(conn, "xqb_users").
			FromSubquery(sub, "a").
			Select("name").
			OrderBy("name", "ASC").
			Get()
		require.NoError(t, err)
		require.Len(t, rows, 2)
	})
}

// TestRawQuery checks raw SQL still flows through the same connection.
func TestRawQuery(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		// Build placeholder-correct SQL via the builder, then run it raw.
		qb := QB(conn, "xqb_users").Select("name").Where("status", "=", "active").OrderBy("name", "ASC")
		sqlStr, args, err := qb.ToSql()
		require.NoError(t, err)

		rows, err := xqb.Sql(sqlStr, args...).Connection(conn).Query()
		require.NoError(t, err)
		defer rows.Close()

		count := 0
		for rows.Next() {
			count++
		}
		require.NoError(t, rows.Err())
		assert.Equal(t, 2, count)
	})
}

func TestModelBind(t *testing.T) {
	type User struct {
		ID     int64  `xqb:"id"`
		Name   string `xqb:"name"`
		Email  string `xqb:"email"`
		Status string `xqb:"status"`
	}
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		data, err := QB(conn, "xqb_users").
			Select("id", "name", "email", "status").
			Where("name", "=", "Alice").
			Get()
		require.NoError(t, err)
		require.Len(t, data, 1)

		var u User
		require.NoError(t, xqb.Bind(data[0], &u))
		assert.Equal(t, "Alice", u.Name)
		assert.Equal(t, "alice@example.com", u.Email)

		all, err := QB(conn, "xqb_users").
			Select("id", "name", "email", "status").
			OrderBy("id", "ASC").
			Get()
		require.NoError(t, err)
		var users []User
		require.NoError(t, xqb.Bind(all, &users))
		require.Len(t, users, 4)
		assert.Equal(t, "Bob", users[1].Name)
	})
}
