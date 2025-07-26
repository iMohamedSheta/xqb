package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/stretchr/testify/assert"
)

func Test_Update(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		data := map[string]any{
			"first_name": "John",
			"last_name":  "Doe",
			"email":      "john@example",
		}

		sql, bindings, err := qb.UpdateSql(data)

		assert.ErrorIs(t, err, xqbErr.ErrInvalidQuery) // update without where clause is dangerous [not allowed]
		assert.Empty(t, sql)
		assert.Empty(t, bindings)
	})
}

func Test_UpdateWhere(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		data := map[string]any{
			"first_name": "John",
			"last_name":  "Doe",
			"email":      "john@example",
		}

		sql, bindings, err := qb.Where("id", "=", 1).UpdateSql(data)

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "UPDATE `users` SET `email` = ?, `first_name` = ?, `last_name` = ? WHERE `id` = ?",
			types.DriverPostgres: `UPDATE "users" SET "email" = $1, "first_name" = $2, "last_name" = $3 WHERE "id" = $4`,
		}

		assert.NoError(t, err)
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"john@example", "John", "Doe", 1}, bindings)
	})
}

func Test_Update_AllowDangerous(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).AllowDangerous()
		data := map[string]any{
			"first_name": "John",
			"last_name":  "Doe",
			"email":      "john@example",
		}

		sql, bindings, err := qb.UpdateSql(data)

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "UPDATE `users` SET `email` = ?, `first_name` = ?, `last_name` = ?",
			types.DriverPostgres: `UPDATE "users" SET "email" = $1, "first_name" = $2, "last_name" = $3`,
		}

		assert.NoError(t, err)
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"john@example", "John", "Doe"}, bindings)
	})
}

func Test_UpdateWithExpressionValue(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).
			Where("id", "=", 1)

		data := map[string]any{
			"login_count": xqb.Raw("login_count + 1"),
		}

		sql, bindings, err := qb.UpdateSql(data)

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "UPDATE `users` SET `login_count` = login_count + 1 WHERE `id` = ?",
			types.DriverPostgres: `UPDATE "users" SET "login_count" = login_count + 1 WHERE "id" = $1`,
		}

		assert.NoError(t, err)
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{1}, bindings)
	})
}
