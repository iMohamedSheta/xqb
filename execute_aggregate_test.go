package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/stretchr/testify/assert"
)

func Test_CountSql(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.CountSql("id")

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT COUNT(`id`) AS `count` FROM `users`",
			types.DriverPostgres: `SELECT COUNT("id") AS "count" FROM "users"`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_AvgSql(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.AvgSql("age")

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT AVG(`age`) AS `avg` FROM `users`",
			types.DriverPostgres: `SELECT AVG("age") AS "avg" FROM "users"`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_SumSql(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.SumSql("points")

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT SUM(`points`) AS `sum` FROM `users`",
			types.DriverPostgres: `SELECT SUM("points") AS "sum" FROM "users"`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_MinSql(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.MinSql("salary")

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT MIN(`salary`) AS `min` FROM `users`",
			types.DriverPostgres: `SELECT MIN("salary") AS "min" FROM "users"`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_MaxSql(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.MaxSql("score")

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT MAX(`score`) AS `max` FROM `users`",
			types.DriverPostgres: `SELECT MAX("score") AS "max" FROM "users"`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_Count_With_Conditions(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).
			Where("status", "=", "active").
			WhereGroup(func(q *xqb.QueryBuilder) {
				q.Where("role", "=", "admin").
					OrWhere("created_at", ">", "2024-01-01")
			}).
			WhereRaw("deleted_at IS NULL")

		sql, bindings, err := qb.CountSql("id")

		expected := map[types.Driver]string{
			types.DriverMySql: "SELECT COUNT(`id`) AS `count` FROM `users` WHERE " +
				"`status` = ? AND (`role` = ? OR `created_at` > ?) AND deleted_at IS NULL",
			types.DriverPostgres: `SELECT COUNT("id") AS "count" FROM "users" WHERE ` +
				`"status" = $1 AND ("role" = $2 OR "created_at" > $3) AND deleted_at IS NULL`,
		}

		assert.Equal(t, expected[dialect], sql)
		assert.Equal(t, []any{"active", "admin", "2024-01-01"}, bindings)
		assert.NoError(t, err)
	})
}
