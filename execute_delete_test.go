package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/stretchr/testify/assert"
)

func Test_Delete(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect)

		sql, bindings, err := qb.DeleteSql()

		assert.ErrorIs(t, err, xqbErr.ErrInvalidQuery) // delete without where clause is dangerous [not allowed]
		assert.Empty(t, sql)
		assert.Empty(t, bindings)
	})
}

func Test_DeleteWhere(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect)

		sql, bindings, err := qb.Where("id", "=", 1).DeleteSql()

		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:     "DELETE FROM `users` WHERE `id` = ?",
			types.DialectMariaDB:   "DELETE FROM `users` WHERE `id` = ?",
			types.DialectPostgres:  `DELETE FROM "users" WHERE "id" = $1`,
			types.DialectSQLite:    `DELETE FROM "users" WHERE "id" = ?`,
			types.DialectSQLServer: `DELETE FROM [users] WHERE [id] = @p1`,
		}

		assert.Equal(t, expectedSQL[dialect], sql)
		assert.Equal(t, []any{1}, bindings)
		assert.NoError(t, err)
	})
}

func Test_DeleteWithLimit(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).
			Where("status", "=", "inactive").
			Limit(10)

		sql, bindings, err := qb.DeleteSql()

		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:     "DELETE FROM `users` WHERE `status` = ? LIMIT 10",
			types.DialectMariaDB:   "DELETE FROM `users` WHERE `status` = ? LIMIT 10",
			types.DialectPostgres:  ``, // PostgreSQL doesn't support LIMIT on DELETE
			types.DialectSQLite:    ``, // SQLite doesn't support LIMIT on DELETE
			types.DialectSQLServer: `DELETE TOP (10) FROM [users] WHERE [status] = @p1`,
		}

		expectedErr := map[types.Dialect]error{
			types.DialectMySql:     nil,
			types.DialectMariaDB:   nil,
			types.DialectPostgres:  xqbErr.ErrInvalidQuery,
			types.DialectSQLite:    xqbErr.ErrInvalidQuery,
			types.DialectSQLServer: nil,
		}

		assert.Equal(t, expectedSQL[dialect], sql)

		if expectedErr[dialect] != nil {
			assert.Empty(t, bindings)
			assert.ErrorIs(t, err, expectedErr[dialect])
		} else {
			assert.NoError(t, err)
			assert.Equal(t, []any{"inactive"}, bindings)
		}
	})
}

func Test_DeleteWithOffset(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).
			Where("status", "=", "inactive").
			Offset(10)

		sql, bindings, err := qb.DeleteSql()

		// Both MySQL and PostgreSQL should return an error — OFFSET not supported in DELETE
		expectedErr := xqbErr.ErrInvalidQuery

		assert.Empty(t, sql)
		assert.Empty(t, bindings)
		assert.ErrorIs(t, err, expectedErr)
	})
}
