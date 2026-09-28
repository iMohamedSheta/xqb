package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/stretchr/testify/assert"
)

func TestLimit(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Select("*").Limit(10)
		sql, bindings, err := qb.ToSql()
		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:     "SELECT * FROM `users` LIMIT 10",
			types.DialectMariaDB:   "SELECT * FROM `users` LIMIT 10",
			types.DialectPostgres:  `SELECT * FROM "users" LIMIT 10`,
			types.DialectSQLite:    `SELECT * FROM "users" LIMIT 10`,
			types.DialectSQLServer: "SELECT TOP 10 * FROM [users]",
		}
		assert.Equal(t, expectedSQL[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func TestOffset(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Select("*").Offset(5)
		sql, bindings, err := qb.ToSql()

		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:     "SELECT * FROM `users` OFFSET 5",
			types.DialectMariaDB:   "SELECT * FROM `users` OFFSET 5",
			types.DialectPostgres:  `SELECT * FROM "users" OFFSET 5`,
			types.DialectSQLite:    `SELECT * FROM "users" OFFSET 5`,
			types.DialectSQLServer: "SELECT * FROM [users] ORDER BY (SELECT 0) OFFSET 5 ROWS",
		}

		assert.Equal(t, expectedSQL[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func TestSkipAlias(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Select("*").Skip(7)
		sql, bindings, err := qb.ToSql()

		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:     "SELECT * FROM `users` OFFSET 7",
			types.DialectMariaDB:   "SELECT * FROM `users` OFFSET 7",
			types.DialectPostgres:  `SELECT * FROM "users" OFFSET 7`,
			types.DialectSQLite:    `SELECT * FROM "users" OFFSET 7`,
			types.DialectSQLServer: "SELECT * FROM [users] ORDER BY (SELECT 0) OFFSET 7 ROWS",
		}

		assert.Equal(t, expectedSQL[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})

}

func TestTakeAlias(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Select("*").Take(25)
		sql, bindings, err := qb.ToSql()

		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:     "SELECT * FROM `users` LIMIT 25",
			types.DialectMariaDB:   "SELECT * FROM `users` LIMIT 25",
			types.DialectPostgres:  `SELECT * FROM "users" LIMIT 25`,
			types.DialectSQLite:    `SELECT * FROM "users" LIMIT 25`,
			types.DialectSQLServer: "SELECT TOP 25 * FROM [users]",
		}

		assert.Equal(t, expectedSQL[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func TestForPage(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Select("*").ForPage(3, 15)
		sql, bindings, err := qb.ToSql()

		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:     "SELECT * FROM `users` LIMIT 15 OFFSET 30",
			types.DialectMariaDB:   "SELECT * FROM `users` LIMIT 15 OFFSET 30",
			types.DialectPostgres:  `SELECT * FROM "users" LIMIT 15 OFFSET 30`,
			types.DialectSQLite:    `SELECT * FROM "users" LIMIT 15 OFFSET 30`,
			types.DialectSQLServer: "SELECT * FROM [users] ORDER BY (SELECT 0) OFFSET 30 ROWS FETCH NEXT 15 ROWS ONLY",
		}

		assert.Equal(t, expectedSQL[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func TestLimitOffsetWithWhere(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("products").SetDialect(dialect).
			Select("id", "name").
			Where("price", ">", 100).
			OrderBy("created_at", "desc").
			Limit(20).
			Offset(40)

		sql, bindings, err := qb.ToSql()

		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:     "SELECT `id`, `name` FROM `products` WHERE `price` > ? ORDER BY `created_at` desc LIMIT 20 OFFSET 40",
			types.DialectMariaDB:   "SELECT `id`, `name` FROM `products` WHERE `price` > ? ORDER BY `created_at` desc LIMIT 20 OFFSET 40",
			types.DialectPostgres:  `SELECT "id", "name" FROM "products" WHERE "price" > $1 ORDER BY "created_at" desc LIMIT 20 OFFSET 40`,
			types.DialectSQLite:    `SELECT "id", "name" FROM "products" WHERE "price" > ? ORDER BY "created_at" desc LIMIT 20 OFFSET 40`,
			types.DialectSQLServer: "SELECT [id], [name] FROM [products] WHERE [price] > @p1 ORDER BY [created_at] desc OFFSET 40 ROWS FETCH NEXT 20 ROWS ONLY",
		}

		assert.Equal(t, expectedSQL[dialect], sql)
		assert.Equal(t, []any{100}, bindings)
		assert.NoError(t, err)
	})
}

func TestForPageWithWhereAndOrder(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("orders").SetDialect(dialect).
			Select("id", "user_id").
			Where("status", "=", "pending").
			OrderBy("id", "ASC").
			ForPage(5, 10)

		sql, bindings, err := qb.ToSql()

		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:     "SELECT `id`, `user_id` FROM `orders` WHERE `status` = ? ORDER BY `id` ASC LIMIT 10 OFFSET 40",
			types.DialectMariaDB:   "SELECT `id`, `user_id` FROM `orders` WHERE `status` = ? ORDER BY `id` ASC LIMIT 10 OFFSET 40",
			types.DialectPostgres:  `SELECT "id", "user_id" FROM "orders" WHERE "status" = $1 ORDER BY "id" ASC LIMIT 10 OFFSET 40`,
			types.DialectSQLite:    `SELECT "id", "user_id" FROM "orders" WHERE "status" = ? ORDER BY "id" ASC LIMIT 10 OFFSET 40`,
			types.DialectSQLServer: "SELECT [id], [user_id] FROM [orders] WHERE [status] = @p1 ORDER BY [id] ASC OFFSET 40 ROWS FETCH NEXT 10 ROWS ONLY",
		}

		assert.Equal(t, expectedSQL[dialect], sql)
		assert.Equal(t, []any{"pending"}, bindings)
		assert.NoError(t, err)
	})
}

func TestPaginationWithJoins(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).
			Select("users.id", "profiles.bio").
			Join("profiles", "profiles.user_id = users.id").
			OrderBy("users.created_at", "desc").
			Limit(50).
			Offset(100)

		sql, bindings, err := qb.ToSql()

		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:     "SELECT `users`.`id`, `profiles`.`bio` FROM `users` JOIN `profiles` ON profiles.user_id = users.id ORDER BY `users`.`created_at` desc LIMIT 50 OFFSET 100",
			types.DialectMariaDB:   "SELECT `users`.`id`, `profiles`.`bio` FROM `users` JOIN `profiles` ON profiles.user_id = users.id ORDER BY `users`.`created_at` desc LIMIT 50 OFFSET 100",
			types.DialectPostgres:  `SELECT "users"."id", "profiles"."bio" FROM "users" JOIN "profiles" ON profiles.user_id = users.id ORDER BY "users"."created_at" desc LIMIT 50 OFFSET 100`,
			types.DialectSQLite:    `SELECT "users"."id", "profiles"."bio" FROM "users" JOIN "profiles" ON profiles.user_id = users.id ORDER BY "users"."created_at" desc LIMIT 50 OFFSET 100`,
			types.DialectSQLServer: "SELECT [users].[id], [profiles].[bio] FROM [users] JOIN [profiles] ON profiles.user_id = users.id ORDER BY [users].[created_at] desc OFFSET 100 ROWS FETCH NEXT 50 ROWS ONLY",
		}

		assert.Equal(t, expectedSQL[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func TestForPageLargePageNumber(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("logs").SetDialect(dialect).
			Select("*").
			ForPage(999, 1000)

		sql, bindings, err := qb.ToSql()

		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:     "SELECT * FROM `logs` LIMIT 1000 OFFSET 998000",
			types.DialectMariaDB:   "SELECT * FROM `logs` LIMIT 1000 OFFSET 998000",
			types.DialectPostgres:  `SELECT * FROM "logs" LIMIT 1000 OFFSET 998000`,
			types.DialectSQLite:    `SELECT * FROM "logs" LIMIT 1000 OFFSET 998000`,
			types.DialectSQLServer: "SELECT * FROM [logs] ORDER BY (SELECT 0) OFFSET 998000 ROWS FETCH NEXT 1000 ROWS ONLY",
		}

		assert.Equal(t, expectedSQL[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func TestForPageWithGroupByHaving(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("transactions").SetDialect(dialect).
			Select("user_id", xqb.Raw("SUM(amount) as total")).
			GroupBy("user_id").
			Having("SUM(amount)", ">", 1000).
			ForPage(2, 25)

		sql, bindings, err := qb.ToSql()

		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:     "SELECT `user_id`, SUM(amount) as total FROM `transactions` GROUP BY `user_id` HAVING SUM(amount) > ? LIMIT 25 OFFSET 25",
			types.DialectMariaDB:   "SELECT `user_id`, SUM(amount) as total FROM `transactions` GROUP BY `user_id` HAVING SUM(amount) > ? LIMIT 25 OFFSET 25",
			types.DialectPostgres:  `SELECT "user_id", SUM(amount) as total FROM "transactions" GROUP BY "user_id" HAVING SUM(amount) > $1 LIMIT 25 OFFSET 25`,
			types.DialectSQLite:    `SELECT "user_id", SUM(amount) as total FROM "transactions" GROUP BY "user_id" HAVING SUM(amount) > ? LIMIT 25 OFFSET 25`,
			types.DialectSQLServer: "SELECT [user_id], SUM(amount) as total FROM [transactions] GROUP BY [user_id] HAVING SUM(amount) > @p1 ORDER BY (SELECT 0) OFFSET 25 ROWS FETCH NEXT 25 ROWS ONLY",
		}

		assert.Equal(t, expectedSQL[dialect], sql)
		assert.Equal(t, []any{1000}, bindings)
		assert.NoError(t, err)
	})
}
