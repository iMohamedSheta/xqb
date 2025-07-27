package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/stretchr/testify/assert"
)

func Test_Pluck_ValueAndKey(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Where("name", "LIKE", "%mohamed%")

		sql, bindings, err := qb.PluckSql("name", "id")
		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "SELECT `name`, `id` FROM `users` WHERE `name` LIKE ?",
			types.DialectPostgres: `SELECT "name", "id" FROM "users" WHERE "name" LIKE $1`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"%mohamed%"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_Pluck_Value(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Where("name", "LIKE", "%mohamed%")
		sql, bindings, err := qb.PluckSql("name", "")

		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "SELECT `name` FROM `users` WHERE `name` LIKE ?",
			types.DialectPostgres: `SELECT "name" FROM "users" WHERE "name" LIKE $1`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"%mohamed%"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_Pluck_Key(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Where("name", "LIKE", "%mohamed%")
		sql, bindings, err := qb.PluckSql("", "id")

		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "SELECT `id` FROM `users` WHERE `name` LIKE ?",
			types.DialectPostgres: `SELECT "id" FROM "users" WHERE "name" LIKE $1`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"%mohamed%"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_Pluck_NoData(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.PluckSql("name", "id")

		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "SELECT `name`, `id` FROM `users`",
			types.DialectPostgres: `SELECT "name", "id" FROM "users"`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_Pluck_NoData_NoKey_ReturnError(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect)
		qb.Where("id", "=", 1)
		sql, bindings, err := qb.PluckSql("", "")

		assert.Error(t, err)
		assert.ErrorIs(t, err, xqbErr.ErrInvalidQuery)
		assert.Equal(t, 0, len(bindings))
		assert.Equal(t, "", sql)
	})
}

func Test_Pluck_ComplexQuery(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Where("name", "LIKE", "%mohamed%")
		qb.Select("id", "name").Where("id", "=", 1)
		sql, bindings, err := qb.PluckSql("", "")

		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "SELECT `id`, `name` FROM `users` WHERE `name` LIKE ? AND `id` = ?",
			types.DialectPostgres: `SELECT "id", "name" FROM "users" WHERE "name" LIKE $1 AND "id" = $2`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"%mohamed%", 1}, bindings)
		assert.NoError(t, err)
	})
}
