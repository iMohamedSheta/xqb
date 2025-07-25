package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/stretchr/testify/assert"
)

func TestLimit(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).Select("*").Limit(10)
		sql, bindings, err := qb.ToSQL()
		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM `users` LIMIT 10",
			types.DriverPostgres: `SELECT * FROM "users" LIMIT 10`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func TestOffset(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).Select("*").Offset(5)
		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM `users` OFFSET 5",
			types.DriverPostgres: `SELECT * FROM "users" OFFSET 5`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func TestSkipAlias(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).Select("*").Skip(7)
		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM `users` OFFSET 7",
			types.DriverPostgres: `SELECT * FROM "users" OFFSET 7`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})

}

func TestTakeAlias(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).Select("*").Take(25)
		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM `users` LIMIT 25",
			types.DriverPostgres: `SELECT * FROM "users" LIMIT 25`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func TestForPage(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).Select("*").ForPage(3, 15)
		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM `users` LIMIT 15 OFFSET 30",
			types.DriverPostgres: `SELECT * FROM "users" LIMIT 15 OFFSET 30`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func TestLimitOffsetWithWhere(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("products").SetDialect(dialect).
			Select("id", "name").
			Where("price", ">", 100).
			OrderBy("created_at", "desc").
			Limit(20).
			Offset(40)

		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT `id`, `name` FROM `products` WHERE `price` > ? ORDER BY `created_at` desc LIMIT 20 OFFSET 40",
			types.DriverPostgres: `SELECT "id", "name" FROM "products" WHERE "price" > $1 ORDER BY "created_at" desc LIMIT 20 OFFSET 40`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{100}, bindings)
		assert.NoError(t, err)
	})
}

func TestForPageWithWhereAndOrder(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("orders").SetDialect(dialect).
			Select("id", "user_id").
			Where("status", "=", "pending").
			OrderBy("id", "ASC").
			ForPage(5, 10)

		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT `id`, `user_id` FROM `orders` WHERE `status` = ? ORDER BY `id` ASC LIMIT 10 OFFSET 40",
			types.DriverPostgres: `SELECT "id", "user_id" FROM "orders" WHERE "status" = $1 ORDER BY "id" ASC LIMIT 10 OFFSET 40`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"pending"}, bindings)
		assert.NoError(t, err)
	})
}

func TestPaginationWithJoins(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).
			Select("users.id", "profiles.bio").
			Join("profiles", "profiles.user_id = users.id").
			OrderBy("users.created_at", "desc").
			Limit(50).
			Offset(100)

		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT `users`.`id`, `profiles`.`bio` FROM `users` JOIN `profiles` ON profiles.user_id = users.id ORDER BY `users`.`created_at` desc LIMIT 50 OFFSET 100",
			types.DriverPostgres: `SELECT "users"."id", "profiles"."bio" FROM "users" JOIN "profiles" ON profiles.user_id = users.id ORDER BY "users"."created_at" desc LIMIT 50 OFFSET 100`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func TestForPageLargePageNumber(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("logs").SetDialect(dialect).
			Select("*").
			ForPage(999, 1000)

		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM `logs` LIMIT 1000 OFFSET 998000",
			types.DriverPostgres: `SELECT * FROM "logs" LIMIT 1000 OFFSET 998000`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func TestForPageWithGroupByHaving(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("transactions").SetDialect(dialect).
			Select("user_id", xqb.Raw("SUM(amount) as total")).
			GroupBy("user_id").
			Having("SUM(amount)", ">", 1000).
			ForPage(2, 25)

		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT `user_id`, SUM(amount) as total FROM `transactions` GROUP BY `user_id` HAVING SUM(amount) > ? LIMIT 25 OFFSET 25",
			types.DriverPostgres: `SELECT "user_id", SUM(amount) as total FROM "transactions" GROUP BY "user_id" HAVING SUM(amount) > $1 LIMIT 25 OFFSET 25`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{1000}, bindings)
		assert.NoError(t, err)
	})
}
