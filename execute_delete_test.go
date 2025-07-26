package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/stretchr/testify/assert"
)

func Test_Delete(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)

		sql, bindings, err := qb.DeleteSql()

		assert.ErrorIs(t, err, xqbErr.ErrInvalidQuery) // delete without where clause is dangerous [not allowed]
		assert.Empty(t, sql)
		assert.Empty(t, bindings)
	})
}

func Test_DeleteWhere(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)

		sql, bindings, err := qb.Where("id", "=", 1).DeleteSql()

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "DELETE FROM `users` WHERE `id` = ?",
			types.DriverPostgres: `DELETE FROM "users" WHERE "id" = $1`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []interface{}{1}, bindings)
		assert.NoError(t, err)
	})
}

func Test_DeleteWithLimit(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).
			Where("status", "=", "inactive").
			Limit(10)

		sql, bindings, err := qb.DeleteSql()

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "DELETE FROM `users` WHERE `status` = ? LIMIT 10",
			types.DriverPostgres: ``, // PostgreSQL doesn't support LIMIT on DELETE
		}

		expectedErr := map[types.Driver]error{
			types.DriverMySql:    nil,
			types.DriverPostgres: xqbErr.ErrInvalidQuery,
		}

		assert.Equal(t, expectedSql[dialect], sql)

		if expectedErr[dialect] != nil {
			assert.Empty(t, bindings)
			assert.ErrorIs(t, err, expectedErr[dialect])
		} else {
			assert.NoError(t, err)
			assert.Equal(t, []interface{}{"inactive"}, bindings)
		}
	})
}

func Test_DeleteWithOffset(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
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
