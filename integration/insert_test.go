package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInsert_SingleRow insert(): one row, then verify via Get.
func TestInsert_SingleRow(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetUsersEmpty(t, conn)

		err := QB(conn, "xqb_users").Insert([]map[string]any{
			{"name": "Eve", "email": "eve@example.com", "age": 28, "score": 88.5, "status": "active"},
		})
		require.NoError(t, err, "[%s] insert single row", conn)

		assert.Equal(t, int64(1), countUsers(t, conn))

		row, err := QB(conn, "xqb_users").Where("email", "=", "eve@example.com").First()
		require.NoError(t, err)
		assert.Equal(t, "Eve", asString(row["name"]))
		assert.Equal(t, int64(28), asInt64(row["age"]))
	})
}

// TestInsert_MultipleRows batch insert.
func TestInsert_MultipleRows(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetUsersEmpty(t, conn)

		err := QB(conn, "xqb_users").Insert([]map[string]any{
			{"name": "Eve", "email": "eve@example.com", "age": 28},
			{"name": "Frank", "email": "frank@example.com", "age": 33},
			{"name": "Grace", "email": "grace@example.com", "age": nil},
		})
		require.NoError(t, err, "[%s] batch insert", conn)
		assert.Equal(t, int64(3), countUsers(t, conn))
	})
}

// TestInsert_NullableColumns ensures NULL bindings survive on every dialect.
func TestInsert_NullableColumns(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetUsersEmpty(t, conn)

		err := QB(conn, "xqb_users").Insert([]map[string]any{
			{"name": "Null Guy", "email": nil, "age": nil, "score": nil, "status": nil},
		})
		require.NoError(t, err, "[%s] insert nulls", conn)

		row, err := QB(conn, "xqb_users").Where("name", "=", "Null Guy").First()
		require.NoError(t, err)
		assert.Nil(t, row["email"])
		assert.Nil(t, row["age"])
	})
}

// TestInsert_InvalidColumn must surface a DB error on every dialect.
func TestInsert_InvalidColumn(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetUsersEmpty(t, conn)

		err := QB(conn, "xqb_users").Insert([]map[string]any{
			{"no_such_column": "boom"},
		})
		assert.Error(t, err, "[%s] expected error for unknown column", conn)
	})
}

// TestInsertGetId insertGetId.
func TestInsertGetId(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetUsersEmpty(t, conn)

		id, err := QB(conn, "xqb_users").InsertGetId([]map[string]any{
			{"name": "Hank", "email": "hank@example.com", "age": 40},
		})
		require.NoError(t, err, "[%s] insertGetId", conn)
		assert.Greater(t, id, int64(0), "[%s] expected positive id", conn)

		row, err := QB(conn, "xqb_users").Find(id)
		require.NoError(t, err)
		assert.Equal(t, "Hank", asString(row["name"]))
	})
}

// TestUpsert upsert(): insert, then conflict-update.
func TestUpsert(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		// SQL Server has no ON CONFLICT support in xqb (returns ErrUnsupportedFeature).
		if conn == "sqlserver" {
			t.Skip("upsert not supported on sqlserver (needs MERGE)")
		}
		resetUsersEmpty(t, conn)

		affected, err := QB(conn, "xqb_users").Upsert(
			[]map[string]any{
				{"name": "Ivy", "email": "ivy@example.com", "age": 22, "status": "active"},
			},
			[]string{"email"},
			[]string{"age", "status"},
		)
		require.NoError(t, err, "[%s] upsert insert", conn)
		assert.GreaterOrEqual(t, affected, int64(1))

		// Second call with same unique key must UPDATE instead of duplicating.
		affected, err = QB(conn, "xqb_users").Upsert(
			[]map[string]any{
				{"name": "Ivy Renamed", "email": "ivy@example.com", "age": 23, "status": "inactive"},
			},
			[]string{"email"},
			[]string{"age", "status"},
		)
		require.NoError(t, err, "[%s] upsert conflict-update", conn)
		assert.GreaterOrEqual(t, affected, int64(1))

		assert.Equal(t, int64(1), countUsers(t, conn), "[%s] upsert must not duplicate", conn)

		row, err := QB(conn, "xqb_users").Where("email", "=", "ivy@example.com").First()
		require.NoError(t, err)
		assert.Equal(t, int64(23), asInt64(row["age"]))
		assert.Equal(t, "inactive", asString(row["status"]))
	})
}
