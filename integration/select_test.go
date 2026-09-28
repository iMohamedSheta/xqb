package integration

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSelect_GetAll get(): full table scan.
func TestSelect_GetAll(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		rows, err := QB(conn, "xqb_users").Select("name", "email").OrderBy("id", "ASC").Get()
		require.NoError(t, err)
		require.Len(t, rows, 4)
		assert.Equal(t, "Alice", asString(rows[0]["name"]))
		assert.Equal(t, "Bob", asString(rows[1]["name"]))
	})
}

// TestSelect_WhereBasic where().
func TestSelect_WhereBasic(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		rows, err := QB(conn, "xqb_users").
			Select("name").
			Where("status", "=", "active").
			OrderBy("name", "ASC").
			Get()
		require.NoError(t, err)
		require.Len(t, rows, 2)
		assert.Equal(t, "Alice", asString(rows[0]["name"]))
		assert.Equal(t, "Bob", asString(rows[1]["name"]))
	})
}

// TestSelect_First first().
func TestSelect_First(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		row, err := QB(conn, "xqb_users").
			Select("name", "email").
			Where("name", "=", "Carol").
			First()
		require.NoError(t, err)
		assert.Equal(t, "Carol", asString(row["name"]))
		assert.Equal(t, "carol@example.com", asString(row["email"]))
	})
}

// TestSelect_Find find(id).
func TestSelect_Find(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		first, err := QB(conn, "xqb_users").OrderBy("id", "ASC").First()
		require.NoError(t, err)
		id := first["id"]

		found, err := QB(conn, "xqb_users").Find(id)
		require.NoError(t, err)
		assert.Equal(t, asString(first["name"]), asString(found["name"]))
	})
}

// TestSelect_FindOrFail findOrFail().
func TestSelect_FindOrFail(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		first, err := QB(conn, "xqb_users").OrderBy("id", "ASC").First()
		require.NoError(t, err)

		row, err := QB(conn, "xqb_users").FindOrFail(first["id"])
		require.NoError(t, err)
		assert.Equal(t, asString(first["name"]), asString(row["name"]))

		_, err = QB(conn, "xqb_users").FindOrFail(-999999)
		assert.Error(t, err, "[%s] findOrFail missing id must error", conn)
		assert.ErrorIs(t, err, xqb.ErrNotFound)
	})
}

// TestSelect_Value value().
func TestSelect_Value(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		v, err := QB(conn, "xqb_users").Where("name", "=", "Alice").Value("email")
		require.NoError(t, err)
		assert.Equal(t, "alice@example.com", asString(v))
	})
}

// TestSelect_Pluck pluck().
func TestSelect_Pluck(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		vals, err := QB(conn, "xqb_users").OrderBy("id", "ASC").PluckSlice("name")
		require.NoError(t, err)
		require.Len(t, vals, 4)
		assert.Equal(t, "Alice", asString(vals[0]))

		mapped, err := QB(conn, "xqb_users").OrderBy("id", "ASC").PluckMap("email", "name")
		require.NoError(t, err)
		assert.Equal(t, "alice@example.com", asString(mapped["Alice"]))
		assert.Equal(t, "bob@example.com", asString(mapped["Bob"]))
	})
}

// TestSelect_Exists exists()/doesntExist().
func TestSelect_Exists(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		ok, err := QB(conn, "xqb_users").Where("email", "=", "alice@example.com").Exists()
		require.NoError(t, err)
		assert.True(t, ok)

		ok, err = QB(conn, "xqb_users").Where("email", "=", "nobody@example.com").Exists()
		require.NoError(t, err)
		assert.False(t, ok)

		missing, err := QB(conn, "xqb_users").Where("email", "=", "nobody@example.com").DoesntExist()
		require.NoError(t, err)
		assert.True(t, missing)
	})
}

// TestSelect_DistinctOrderLimit distinct/orderBy/limit/offset.
func TestSelect_DistinctOrderLimit(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		rows, err := QB(conn, "xqb_users").
			Select("status").
			Distinct().
			WhereNotNull("status").
			OrderBy("status", "ASC").
			Get()
		require.NoError(t, err)
		require.Len(t, rows, 2, "[%s] distinct statuses", conn)

		rows, err = QB(conn, "xqb_users").
			Select("name").
			WhereNotNull("age").
			OrderBy("age", "DESC").
			Limit(2).
			Get()
		require.NoError(t, err)
		require.Len(t, rows, 2)
		// Carol (35) then Alice (30); NULL ages are excluded because
		// NULL ordering under ORDER BY .. DESC is dialect-specific
		// (postgres sorts NULLS FIRST, others sort them last).
		assert.Equal(t, "Carol", asString(rows[0]["name"]))

		rows, err = QB(conn, "xqb_users").
			Select("name").
			OrderBy("id", "ASC").
			Limit(1).Offset(1).
			Get()
		require.NoError(t, err)
		require.Len(t, rows, 1)
		assert.Equal(t, "Bob", asString(rows[0]["name"]))
	})
}

// TestSelect_Paginate paginate().
func TestSelect_Paginate(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		rows, meta, err := QB(conn, "xqb_users").
			Select("id", "name").
			OrderBy("id", "ASC").
			Paginate(2, 1, "*")
		require.NoError(t, err)
		require.Len(t, rows, 2)
		assert.Equal(t, 2, meta["per_page"])
		assert.Equal(t, 1, meta["current_page"])
		assert.Equal(t, int64(4), asInt64(meta["total_count"]))

		rows, meta, err = QB(conn, "xqb_users").
			Select("id", "name").
			OrderBy("id", "ASC").
			Paginate(2, 2, "*")
		require.NoError(t, err)
		require.Len(t, rows, 2)
		assert.Equal(t, "Carol", asString(rows[0]["name"]))
		assert.Equal(t, 2, meta["current_page"])
	})
}

// TestSelect_Chunks chunk().
func TestSelect_Chunks(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		total := 0
		chunks := 0
		err := QB(conn, "xqb_users").
			Select("id", "name").
			OrderBy("id", "ASC").
			Chunks(2, func(rows []map[string]any) error {
				chunks++
				total += len(rows)
				return nil
			})
		require.NoError(t, err)
		assert.Equal(t, 4, total)
		assert.Equal(t, 2, chunks)
	})
}
