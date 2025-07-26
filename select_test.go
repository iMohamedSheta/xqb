package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/stretchr/testify/assert"
)

func Test_Select(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		qb.Select("id", "name", "email")
		sql, bindings, err := qb.ToSql()

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT `id`, `name`, `email` FROM `users`",
			types.DriverPostgres: `SELECT "id", "name", "email" FROM "users"`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_Select_WithWhere(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		qb.Select("id", "name")
		qb.Where("age", ">", 18)
		sql, bindings, err := qb.ToSql()
		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT `id`, `name` FROM `users` WHERE `age` > ?",
			types.DriverPostgres: `SELECT "id", "name" FROM "users" WHERE "age" > $1`,
		}
		expectedBindings := []any{18}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, expectedBindings, bindings)
		assert.NoError(t, err)
	})
}

func Test_Select_WithJoins(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		qb.Select("users.id", "users.name", "orders.id as order_id")
		qb.Join("orders", "users.id = orders.user_id")
		sql, bindings, err := qb.ToSql()
		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT `users`.`id`, `users`.`name`, `orders`.`id` AS `order_id` FROM `users` JOIN `orders` ON users.id = orders.user_id",
			types.DriverPostgres: `SELECT "users"."id", "users"."name", "orders"."id" AS "order_id" FROM "users" JOIN "orders" ON users.id = orders.user_id`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_Select_WithLeftJoins(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		qb.Select("users.id", "users.name", "orders.id as order_id").Where("users.id", ">", 55)
		qb.Join("orders", "users.id = orders.user_id").Where("orders.id", ">", 11)
		qb.LeftJoin("products", "orders.product_id = products.id")
		sql, bindings, err := qb.ToSql()
		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT `users`.`id`, `users`.`name`, `orders`.`id` AS `order_id` FROM `users` JOIN `orders` ON users.id = orders.user_id LEFT JOIN `products` ON orders.product_id = products.id WHERE `users`.`id` > ? AND `orders`.`id` > ?",
			types.DriverPostgres: `SELECT "users"."id", "users"."name", "orders"."id" AS "order_id" FROM "users" JOIN "orders" ON users.id = orders.user_id LEFT JOIN "products" ON orders.product_id = products.id WHERE "users"."id" > $1 AND "orders"."id" > $2`,
		}
		expectedBindings := []any{55, 11}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, expectedBindings, bindings)
		assert.NoError(t, err)
	})
}

func Test_Select_WithGroupBy(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("orders").SetDialect(dialect)
		qb.Select("user_id", "COUNT(*) as order_count")
		qb.GroupBy("user_id")
		sql, bindings, err := qb.ToSql()
		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT `user_id`, COUNT(*) AS `order_count` FROM `orders` GROUP BY `user_id`",
			types.DriverPostgres: `SELECT "user_id", COUNT(*) AS "order_count" FROM "orders" GROUP BY "user_id"`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_Select_WithHaving(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("orders").SetDialect(dialect)
		qb.Select("user_id", "COUNT(*) as order_count")
		qb.GroupBy("user_id")
		qb.Having("order_count", ">", 5)
		sql, bindings, err := qb.ToSql()
		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT `user_id`, COUNT(*) AS `order_count` FROM `orders` GROUP BY `user_id` HAVING `order_count` > ?",
			types.DriverPostgres: `SELECT "user_id", COUNT(*) AS "order_count" FROM "orders" GROUP BY "user_id" HAVING "order_count" > $1`,
		}
		expectedBindings := []any{5}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, expectedBindings, bindings)
		assert.NoError(t, err)
	})
}

func Test_Select_WithOrderBy(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		qb.Select("id", "name")
		qb.OrderBy("name", "ASC")
		sql, bindings, err := qb.ToSql()
		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT `id`, `name` FROM `users` ORDER BY `name` ASC",
			types.DriverPostgres: `SELECT "id", "name" FROM "users" ORDER BY "name" ASC`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_Select_WithLimitOffset(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		qb.Select("id", "name")
		qb.Limit(10)
		qb.Offset(20)
		sql, bindings, err := qb.ToSql()
		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT `id`, `name` FROM `users` LIMIT 10 OFFSET 20",
			types.DriverPostgres: `SELECT "id", "name" FROM "users" LIMIT 10 OFFSET 20`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_Select_WithAggregateFunctions(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("orders").SetDialect(dialect).
			Select(
				xqb.Sum("amount", "total_amount"),
				xqb.Avg("amount", "average_amount"),
				xqb.Count("id", "order_count"),
			)

		sql, bindings, err := qb.ToSql()

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT SUM(amount) AS total_amount, AVG(amount) AS average_amount, COUNT(id) AS order_count FROM `orders`",
			types.DriverPostgres: `SELECT SUM(amount) AS total_amount, AVG(amount) AS average_amount, COUNT(id) AS order_count FROM "orders"`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_Select_WithCTE(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		qb.WithRaw("user_totals", "SELECT user_id, SUM(amount) as total_spent FROM orders GROUP BY user_id")
		qb.Select("users.id", "users.name", "user_totals.total_spent")
		qb.Join("user_totals", "users.id = user_totals.user_id")
		sql, bindings, err := qb.ToSql()

		expectedSql := map[types.Driver]string{
			types.DriverMySql: "WITH user_totals AS (SELECT user_id, SUM(amount) as total_spent FROM orders GROUP BY user_id) " +
				"SELECT `users`.`id`, `users`.`name`, `user_totals`.`total_spent` FROM `users` JOIN `user_totals` ON users.id = user_totals.user_id",
			types.DriverPostgres: `WITH user_totals AS (SELECT user_id, SUM(amount) as total_spent FROM orders GROUP BY user_id) ` +
				`SELECT "users"."id", "users"."name", "user_totals"."total_spent" FROM "users" JOIN "user_totals" ON users.id = user_totals.user_id`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_Select_WithComplexCTE(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("products").SetDialect(dialect)
		qb.WithRaw("active_users",
			"WITH user_orders AS (SELECT user_id, COUNT(*) as order_count FROM orders GROUP BY user_id) "+
				"SELECT users.id, users.name, user_orders.order_count FROM users "+
				"JOIN user_orders ON users.id = user_orders.user_id")
		qb.Select("products.id", "products.name", "active_users.name as buyer")
		qb.Join("active_users", "products.id = active_users.id")
		sql, bindings, err := qb.ToSql()

		expectedSql := map[types.Driver]string{
			types.DriverMySql: "WITH active_users AS (WITH user_orders AS (SELECT user_id, COUNT(*) as order_count FROM orders GROUP BY user_id) " +
				"SELECT users.id, users.name, user_orders.order_count FROM users JOIN user_orders ON users.id = user_orders.user_id) " +
				"SELECT `products`.`id`, `products`.`name`, `active_users`.`name` AS `buyer` FROM `products` JOIN `active_users` ON products.id = active_users.id",
			types.DriverPostgres: `WITH active_users AS (WITH user_orders AS (SELECT user_id, COUNT(*) as order_count FROM orders GROUP BY user_id) ` +
				`SELECT users.id, users.name, user_orders.order_count FROM users JOIN user_orders ON users.id = user_orders.user_id) ` +
				`SELECT "products"."id", "products"."name", "active_users"."name" AS "buyer" FROM "products" JOIN "active_users" ON products.id = active_users.id`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_Select_WithJSONExpressions(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		qb.Select(
			"id",
			"name",
			xqb.JsonExtract("metadata", "preferences.theme", "theme"),
		)
		sql, bindings, err := qb.ToSql()

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT `id`, `name`, JSON_EXTRACT(metadata, '$.preferences.theme') AS theme FROM `users`",
			types.DriverPostgres: `SELECT "id", "name", metadata->'preferences'->>'theme' AS theme FROM "users"`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_Select_WithStringFunctions(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		qb.Select(
			"id",
			xqb.Concat([]string{
				"first_name",
				"' '",
				"last_name",
			}, "full_name"),
		)
		sql, bindings, err := qb.ToSql()

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT `id`, CONCAT(first_name, ' ', last_name) AS full_name FROM `users`",
			types.DriverPostgres: `SELECT "id", CONCAT(first_name, ' ', last_name) AS full_name FROM "users"`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any(nil), bindings)
		assert.NoError(t, err)
	})
}

func Test_Select_WithDateFunctions(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("orders").SetDialect(dialect)
		qb.Select(
			"id",
			xqb.DateFormat("created_at", "%Y-%m-%d", "order_date"),
		)
		sql, bindings, err := qb.ToSql()

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT `id`, DATE_FORMAT(created_at, '%Y-%m-%d') AS order_date FROM `orders`",
			types.DriverPostgres: `SELECT "id", TO_CHAR(created_at, '%Y-%m-%d') AS order_date FROM "orders"`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any(nil), bindings)
		assert.NoError(t, err)
	})
}

func Test_Select_WithMathExpressions(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("orders").SetDialect(dialect)
		qb.Select(
			"id",
			xqb.Math("amount * 1.1", "total_with_tax"),
		)
		sql, bindings, err := qb.ToSql()

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT `id`, amount * 1.1 AS total_with_tax FROM `orders`",
			types.DriverPostgres: `SELECT "id", amount * 1.1 AS total_with_tax FROM "orders"`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_Select_WithLocking(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		qb.Select("id", "name")
		qb.LockForUpdate()
		sql, bindings, err := qb.ToSql()

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT `id`, `name` FROM `users` FOR UPDATE",
			types.DriverPostgres: `SELECT "id", "name" FROM "users" FOR UPDATE`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_Select_WithUnion(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).Select("id", "name")
		qb.UnionRaw("SELECT id, name FROM users WHERE type = $1", "admin") // Raw sql will be the same as the raw in any dialect
		qb.UnionRaw("SELECT id, name FROM users WHERE type = $2", "superuser")
		qb.UnionRaw("SELECT id, name FROM users WHERE type = $3", "guest")

		sql, bindings, err := qb.ToSql()

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "(SELECT `id`, `name` FROM `users`) UNION (SELECT id, name FROM users WHERE type = $1) UNION (SELECT id, name FROM users WHERE type = $2) UNION (SELECT id, name FROM users WHERE type = $3)",
			types.DriverPostgres: `(SELECT "id", "name" FROM "users") UNION (SELECT id, name FROM users WHERE type = $1) UNION (SELECT id, name FROM users WHERE type = $2) UNION (SELECT id, name FROM users WHERE type = $3)`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		expectedBindings := []any{"admin", "superuser", "guest"}
		assert.Equal(t, expectedBindings, bindings)
		assert.NoError(t, err)
	})
}

func Test_Select_WithDistinct(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		qb.Select("name")
		qb.Distinct()
		sql, bindings, err := qb.ToSql()

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT DISTINCT `name` FROM `users`",
			types.DriverPostgres: `SELECT DISTINCT "name" FROM "users"`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_Select_WithRawExpressions(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.Select(
			xqb.Raw("COUNT(*) as total"),
			"name",
			xqb.Raw("CONCAT(first_name, ' ', last_name) as full_name"),
		).ToSql()

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT COUNT(*) as total, `name`, CONCAT(first_name, ' ', last_name) as full_name FROM `users`",
			types.DriverPostgres: `SELECT COUNT(*) as total, "name", CONCAT(first_name, ' ', last_name) as full_name FROM "users"`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_Select_WithDateExpressions(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("orders").SetDialect(dialect)
		sql, bindings, err := qb.Select(
			xqb.DateFormat("created_at", "%Y-%m", "month"),
			xqb.Raw("COUNT(*) as total_orders"),
			xqb.Raw("SUM(amount) as total_amount"),
		).
			GroupBy(xqb.DateFormat("created_at", "%Y-%m", "")).
			OrderBy(xqb.DateFormat("created_at", "%Y-%m", ""), "ASC").
			ToSql()

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT DATE_FORMAT(created_at, '%Y-%m') AS month, COUNT(*) as total_orders, SUM(amount) as total_amount FROM `orders` GROUP BY DATE_FORMAT(created_at, '%Y-%m') ORDER BY DATE_FORMAT(created_at, '%Y-%m') ASC",
			types.DriverPostgres: `SELECT TO_CHAR(created_at, '%Y-%m') AS month, COUNT(*) as total_orders, SUM(amount) as total_amount FROM "orders" GROUP BY TO_CHAR(created_at, '%Y-%m') ORDER BY TO_CHAR(created_at, '%Y-%m') ASC`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_Select_WithExpressions(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.Select(
			"id",
			xqb.Raw("CONCAT(first_name, ' ', last_name) as full_name"),
			xqb.Raw("(SELECT COUNT(*) FROM orders WHERE orders.user_id = users.id) as order_count"),
		).
			Where(xqb.Raw("LOWER(email)"), "LIKE", "%@example.com").
			GroupBy("id", "first_name", "last_name").
			Having(xqb.Raw("(SELECT COUNT(*) FROM orders WHERE orders.user_id = users.id)"), ">", 5).
			OrderBy(xqb.Raw("(SELECT SUM(amount) FROM orders WHERE orders.user_id = users.id)"), "DESC").
			ToSql()

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT `id`, CONCAT(first_name, ' ', last_name) as full_name, (SELECT COUNT(*) FROM orders WHERE orders.user_id = users.id) as order_count FROM `users` WHERE LOWER(email) LIKE ? GROUP BY `id`, `first_name`, `last_name` HAVING (SELECT COUNT(*) FROM orders WHERE orders.user_id = users.id) > ? ORDER BY (SELECT SUM(amount) FROM orders WHERE orders.user_id = users.id) DESC",
			types.DriverPostgres: `SELECT "id", CONCAT(first_name, ' ', last_name) as full_name, (SELECT COUNT(*) FROM orders WHERE orders.user_id = users.id) as order_count FROM "users" WHERE LOWER(email) LIKE $1 GROUP BY "id", "first_name", "last_name" HAVING (SELECT COUNT(*) FROM orders WHERE orders.user_id = users.id) > $2 ORDER BY (SELECT SUM(amount) FROM orders WHERE orders.user_id = users.id) DESC`,
		}
		expectedBindings := []any{"%@example.com", 5}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, expectedBindings, bindings)
		assert.NoError(t, err)
	})
}

func Test_Select_WithSubQuery(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		subSql, subBindings, _ := xqb.Table("payments").SetDialect(dialect).
			Select("id", "amount", "created_at").
			Where("payments.user_id", "=", 15).
			ToSql()

		qb := xqb.Table("users").SetDialect(dialect).
			Select("id", "name", xqb.Raw("("+subSql+") AS payments", subBindings...)).
			Where("id", "=", 15)

		sql, bindings, err := qb.ToSql()

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT `id`, `name`, (SELECT `id`, `amount`, `created_at` FROM `payments` WHERE `payments`.`user_id` = ?) AS payments FROM `users` WHERE `id` = ?",
			types.DriverPostgres: `SELECT "id", "name", (SELECT "id", "amount", "created_at" FROM "payments" WHERE "payments"."user_id" = $1) AS payments FROM "users" WHERE "id" = $2`,
		}
		expectedBindings := []any{15, 15}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, expectedBindings, bindings)
		assert.NoError(t, err)
	})
}

func Test_Select_WithSubQuery_(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		sub := xqb.Table("payments").
			Select("id", "amount", "created_at").
			Where("payments.user_id", "=", 15)

		sub2 := xqb.Table("admins").
			Select("id", "amount", "created_at").
			Where("admins.user_id", "=", 15)

		qb := xqb.Table("users").SetDialect(dialect).
			Select("id", "name").
			SelectSub(sub, "payments").
			SelectSub(sub2, "admins").
			Where("id", "=", 15)

		sql, bindings, err := qb.ToSql()

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT `id`, `name`, (SELECT `id`, `amount`, `created_at` FROM `payments` WHERE `payments`.`user_id` = ?) AS payments, (SELECT `id`, `amount`, `created_at` FROM `admins` WHERE `admins`.`user_id` = ?) AS admins FROM `users` WHERE `id` = ?",
			types.DriverPostgres: `SELECT "id", "name", (SELECT "id", "amount", "created_at" FROM "payments" WHERE "payments"."user_id" = $1) AS payments, (SELECT "id", "amount", "created_at" FROM "admins" WHERE "admins"."user_id" = $2) AS admins FROM "users" WHERE "id" = $3`,
		}
		expectedBindings := []any{15, 15, 15}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, expectedBindings, bindings)
		assert.NoError(t, err)
	})
}
func Test_FromSubquery(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		sub := xqb.Table("orders").
			Select("user_id", xqb.Raw("COUNT(*) AS order_count")).
			Where("user_id", "=", 25).
			GroupBy("user_id")

		qb := xqb.New().SetDialect(dialect).
			Select("u.id", "u.name", "o.order_count").
			FromSubquery(sub, "o").
			Join("users u", "u.id = o.user_id").
			Where("u.id", "=", 25)

		sql, bindings, err := qb.ToSql()

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT `u`.`id`, `u`.`name`, `o`.`order_count` FROM (SELECT `user_id`, COUNT(*) AS order_count FROM `orders` WHERE `user_id` = ? GROUP BY `user_id`) AS o JOIN `users` `u` ON u.id = o.user_id WHERE `u`.`id` = ?",
			types.DriverPostgres: `SELECT "u"."id", "u"."name", "o"."order_count" FROM (SELECT "user_id", COUNT(*) AS order_count FROM "orders" WHERE "user_id" = $1 GROUP BY "user_id") AS o JOIN "users" "u" ON u.id = o.user_id WHERE "u"."id" = $2`,
		}
		expectedBindings := []any{25, 25}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, expectedBindings, bindings)
		assert.NoError(t, err)
	})
}
