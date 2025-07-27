package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/stretchr/testify/assert"
)

func Test_CTE_With(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		mainQB := xqb.Query().SetDialect(dialect)
		cteQB := xqb.Table("users").Select("id", "name")
		mainQB.With("cte_users", cteQB)

		cte := mainQB.GetData().WithCTEs[0]

		assert.Len(t, mainQB.GetData().WithCTEs, 1)
		assert.Equal(t, "cte_users", cte.Name)
		assert.NotNil(t, cte.Query)
		assert.Nil(t, cte.Expression)
		assert.False(t, cte.Recursive)

		sql, bindings, err := mainQB.Select("*").ToSql()

		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "WITH cte_users AS (SELECT `id`, `name` FROM `users`) SELECT *",
			types.DialectPostgres: `WITH cte_users AS (SELECT "id", "name" FROM "users") SELECT *`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any(nil), bindings)
		assert.NoError(t, err)
	})
}

func Test_CTE_WithExpression(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		mainQB := xqb.Query().SetDialect(dialect)
		mainQB.WithExpr("cte_expr", "SELECT ?", 42)

		cte := mainQB.GetData().WithCTEs[0]
		assert.Len(t, mainQB.GetData().WithCTEs, 1)
		assert.Equal(t, "cte_expr", cte.Name)
		assert.Nil(t, cte.Query)
		assert.NotNil(t, cte.Expression)
		assert.Equal(t, "SELECT ?", cte.Expression.Sql)
		assert.Equal(t, []any{42}, cte.Expression.Bindings)

		sql, bindings, err := mainQB.Select("*").ToSql()
		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "WITH cte_expr AS (SELECT ?) SELECT *",
			types.DialectPostgres: `WITH cte_expr AS (SELECT $1) SELECT *`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{42}, bindings)
		assert.NoError(t, err)
	})
}

func Test_CTE_WithRecursive(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		mainQB := xqb.Query().SetDialect(dialect)
		recQB := xqb.Table("tree").Select("id", "parent_id")
		mainQB.WithRecursive("cte_tree", recQB)

		cte := mainQB.GetData().WithCTEs[0]
		assert.True(t, cte.Recursive)

		sql, b, err := mainQB.Select("*").ToSql()

		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "WITH RECURSIVE cte_tree AS (SELECT `id`, `parent_id` FROM `tree`) SELECT *",
			types.DialectPostgres: `WITH RECURSIVE cte_tree AS (SELECT "id", "parent_id" FROM "tree") SELECT *`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any(nil), b)
		assert.NoError(t, err)
	})
}

func Test_CTE_WithRaw(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {

		mainQB := xqb.New().SetDialect(dialect)
		mainQB.WithRaw("cte_raw", "SELECT ? AS col", 99)

		cte := mainQB.GetData().WithCTEs[0]
		assert.NotNil(t, cte.Expression)
		assert.False(t, cte.Recursive)

		sql, bindings, err := mainQB.Select("*").ToSql()

		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "WITH cte_raw AS (SELECT ? AS col) SELECT *",
			types.DialectPostgres: `WITH cte_raw AS (SELECT $1 AS col) SELECT *`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{99}, bindings)
		assert.NoError(t, err)
	})
}

func Test_CTE_WithRecursiveRaw(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		mainQB := xqb.Query().SetDialect(dialect)
		mainQB.WithRecursiveRaw("cte_rec_raw", "SELECT ? AS col", 123)

		cte := mainQB.GetData().WithCTEs[0]
		assert.True(t, cte.Recursive)

		sql, bindings, err := mainQB.Select("*").ToSql()

		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "WITH RECURSIVE cte_rec_raw AS (SELECT ? AS col) SELECT *",
			types.DialectPostgres: `WITH RECURSIVE cte_rec_raw AS (SELECT $1 AS col) SELECT *`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{123}, bindings)
		assert.NoError(t, err)
	})
}

func Test_CTE_WithAdvancedExpressions(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		cteQB := xqb.Table("coverage_table").
			Select(
				"status",
				xqb.Sum("amount", "total_amount"),
				xqb.Length("bio", "bio_len"),
			).Where(
			xqb.Lower("status", ""), "=", "active",
		).GroupBy(
			xqb.Date("created_at", ""),
			xqb.Upper("region", ""),
		).Having(
			"total_amount", ">", 1000,
		).OrderBy(
			xqb.Length("bio", ""), "DESC",
		).Limit(5).Offset(10)

		mainQB := xqb.New().SetDialect(dialect)
		mainQB.With("cte_agg", cteQB).Select("*")

		sql, bindings, err := mainQB.ToSql()
		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "WITH cte_agg AS (SELECT `status`, SUM(amount) AS total_amount, LENGTH(bio) AS bio_len FROM `coverage_table` WHERE LOWER(status) = ? GROUP BY DATE(created_at), UPPER(region) HAVING `total_amount` > ? ORDER BY LENGTH(bio) DESC LIMIT 5 OFFSET 10) SELECT *",
			types.DialectPostgres: `WITH cte_agg AS (SELECT "status", SUM(amount) AS total_amount, LENGTH(bio) AS bio_len FROM "coverage_table" WHERE LOWER(status) = $1 GROUP BY DATE(created_at), UPPER(region) HAVING "total_amount" > $2 ORDER BY LENGTH(bio) DESC LIMIT 5 OFFSET 10) SELECT *`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"active", 1000}, bindings)
		assert.NoError(t, err)
	})
}

func Test_CTE_WithMultipleCTEs(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		mainQB := xqb.Query().SetDialect(dialect)
		mainQB.
			WithRaw("cte1", "SELECT 1 AS one").
			WithExpr("cte2", "SELECT 2 AS two").
			With("cte3", xqb.Table("users").Select("id"))

		sql, b, err := mainQB.Select("*").ToSql()
		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "WITH cte1 AS (SELECT 1 AS one), cte2 AS (SELECT 2 AS two), cte3 AS (SELECT `id` FROM `users`) SELECT *",
			types.DialectPostgres: `WITH cte1 AS (SELECT 1 AS one), cte2 AS (SELECT 2 AS two), cte3 AS (SELECT "id" FROM "users") SELECT *`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any(nil), b)
		assert.NoError(t, err)
	})
}

func Test_CTE_WithAliasedExpressions(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		mainQB := xqb.Query().SetDialect(dialect)
		mainQB.WithExpr("agg_stats", "SELECT COUNT(*) AS total, MAX(score) AS high_score FROM games")

		sql, bindings, err := mainQB.Select("*").ToSql()

		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "WITH agg_stats AS (SELECT COUNT(*) AS total, MAX(score) AS high_score FROM games) SELECT *",
			types.DialectPostgres: `WITH agg_stats AS (SELECT COUNT(*) AS total, MAX(score) AS high_score FROM games) SELECT *`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any(nil), bindings)
		assert.NoError(t, err)
	})
}

func Test_CTE_UsageInMainQuery(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		mainQB := xqb.New().SetDialect(dialect)
		mainQB.
			With("cte_users", xqb.New().Table("users").Select("id", "name")).
			From("cte_users").
			Where("id", ">", 5)

		sql, bindings, err := mainQB.ToSql()
		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "WITH cte_users AS (SELECT `id`, `name` FROM `users`) SELECT * FROM `cte_users` WHERE `id` > ?",
			types.DialectPostgres: `WITH cte_users AS (SELECT "id", "name" FROM "users") SELECT * FROM "cte_users" WHERE "id" > $1`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{5}, bindings)
		assert.NoError(t, err)
	})
}

func Test_CTE_Recursive_Usage(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		recQB := xqb.Table("tree").Select("id", "parent_id")
		mainQB := xqb.Table("tree_cte").SetDialect(dialect).
			WithRecursive("tree_cte", recQB).
			WhereNull("parent_id")

		sql, b, err := mainQB.ToSql()
		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "WITH RECURSIVE tree_cte AS (SELECT `id`, `parent_id` FROM `tree`) SELECT * FROM `tree_cte` WHERE `parent_id` IS NULL",
			types.DialectPostgres: `WITH RECURSIVE tree_cte AS (SELECT "id", "parent_id" FROM "tree") SELECT * FROM "tree_cte" WHERE "parent_id" IS NULL`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any(nil), b)
		assert.NoError(t, err)
	})
}

func Test_CTE_BindingsOrder(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		mainQB := xqb.Query().SetDialect(dialect)

		mainQB.
			WithRaw("cte1", "SELECT ? AS one", 1).
			WithRaw("cte2", "SELECT ? AS two", 2).
			From("cte2").
			Where("two", ">", 3)

		sql, bindings, _ := mainQB.ToSql()
		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "WITH cte1 AS (SELECT ? AS one), cte2 AS (SELECT ? AS two) SELECT * FROM `cte2` WHERE `two` > ?",
			types.DialectPostgres: `WITH cte1 AS (SELECT $1 AS one), cte2 AS (SELECT $2 AS two) SELECT * FROM "cte2" WHERE "two" > $3`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{1, 2, 3}, bindings)
	})
}

func Test_CTE_EmptyCTEsShouldNotEmitWith(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Select("id")
		sql, bindings, err := qb.ToSql()
		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "SELECT `id` FROM `users`",
			types.DialectPostgres: `SELECT "id" FROM "users"`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any(nil), bindings)
		assert.NoError(t, err)
	})
}

func Test_CTE_ComplexThreeLevelChain(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		mainQB := xqb.Query().SetDialect(dialect)

		highValueOrders := xqb.Table("orders").Select("user_id", "total").
			Where("total", ">", 100)

		userOrderDetails := xqb.Table("high_value_orders").Select(
			"high_value_orders.user_id",
			"users.name",
		).Join("users", "users.id = high_value_orders.user_id")

		userOrderSummary := xqb.Table("user_order_details").Select(
			"name",
			xqb.Count("*", "order_count"),
		).GroupBy("name")

		mainQB.
			With("high_value_orders", highValueOrders).
			With("user_order_details", userOrderDetails).
			With("user_order_summary", userOrderSummary).
			From("user_order_summary").
			Where("order_count", ">", 5).
			OrderBy("order_count", "DESC")

		sql, bindings, err := mainQB.ToSql()
		expectedSql := map[types.Dialect]string{
			types.DialectMySql: "WITH " +
				"high_value_orders AS (SELECT `user_id`, `total` FROM `orders` WHERE `total` > ?), " +
				"user_order_details AS (SELECT `high_value_orders`.`user_id`, `users`.`name` FROM `high_value_orders` JOIN `users` ON users.id = high_value_orders.user_id), " +
				"user_order_summary AS (SELECT `name`, COUNT(*) AS order_count FROM `user_order_details` GROUP BY `name`) " +
				"SELECT * FROM `user_order_summary` WHERE `order_count` > ? ORDER BY `order_count` DESC",
			types.DialectPostgres: `WITH ` +
				`high_value_orders AS (SELECT "user_id", "total" FROM "orders" WHERE "total" > $1), ` +
				`user_order_details AS (SELECT "high_value_orders"."user_id", "users"."name" FROM "high_value_orders" JOIN "users" ON users.id = high_value_orders.user_id), ` +
				`user_order_summary AS (SELECT "name", COUNT(*) AS order_count FROM "user_order_details" GROUP BY "name") ` +
				`SELECT * FROM "user_order_summary" WHERE "order_count" > $2 ORDER BY "order_count" DESC`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{100, 5}, bindings)
		assert.NoError(t, err)
	})
}
