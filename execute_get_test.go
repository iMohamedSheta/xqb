package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/stretchr/testify/assert"
)

func Test_QueryBuilder_GetSql(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").
			SetDialect(dialect).
			Select("id", "name").
			Where("status", "=", "active").
			OrderBy("created_at", "DESC")

		sql, bindings, err := qb.GetSql()

		expectedSQL := map[types.Driver]string{
			types.DriverMySql:    "SELECT `id`, `name` FROM `users` WHERE `status` = ? ORDER BY `created_at` DESC",
			types.DriverPostgres: `SELECT "id", "name" FROM "users" WHERE "status" = $1 ORDER BY "created_at" DESC`,
		}

		assert.NoError(t, err)
		assert.Equal(t, expectedSQL[dialect], sql)
		assert.Equal(t, []any{"active"}, bindings)
	})
}

func Test_FirstSql(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).Where("id", "=", 10)

		sql, bindings, err := qb.FirstSql()

		expected := map[types.Driver]string{
			types.DriverMySql:    "SELECT * FROM `users` WHERE `id` = ? LIMIT 1",
			types.DriverPostgres: `SELECT * FROM "users" WHERE "id" = $1 LIMIT 1`,
		}

		assert.NoError(t, err)
		assert.Equal(t, expected[dialect], sql)
		assert.Equal(t, []any{10}, bindings)
	})
}

func Test_ValueSql(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).Where("id", "=", 1)

		sql, bindings, err := qb.ValueSql("email")

		expected := map[types.Driver]string{
			types.DriverMySql:    "SELECT `email` FROM `users` WHERE `id` = ? LIMIT 1",
			types.DriverPostgres: `SELECT "email" FROM "users" WHERE "id" = $1 LIMIT 1`,
		}

		assert.NoError(t, err)
		assert.Equal(t, expected[dialect], sql)
		assert.Equal(t, []any{1}, bindings)
	})
}

func Test_FindSql(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)

		sql, bindings, err := qb.FindSql(42)

		expected := map[types.Driver]string{
			types.DriverMySql:    "SELECT * FROM `users` WHERE `id` = ? LIMIT 1",
			types.DriverPostgres: `SELECT * FROM "users" WHERE "id" = $1 LIMIT 1`,
		}

		assert.NoError(t, err)
		assert.Equal(t, expected[dialect], sql)
		assert.Equal(t, []any{42}, bindings)
	})
}

func Test_PaginateSql(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).Where("active", "=", true)

		sql, bindings, err := qb.PaginateSql(10, 3)

		expected := map[types.Driver]string{
			types.DriverMySql:    "SELECT * FROM `users` WHERE `active` = ? LIMIT 10 OFFSET 20",
			types.DriverPostgres: `SELECT * FROM "users" WHERE "active" = $1 LIMIT 10 OFFSET 20`,
		}

		assert.NoError(t, err)
		assert.Equal(t, expected[dialect], sql)
		assert.Equal(t, []any{true}, bindings)
	})
}
