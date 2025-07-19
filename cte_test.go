package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/stretchr/testify/assert"
)

func Test_CTE_With(t *testing.T) {
	mainQB := xqb.Query()
	cteQB := xqb.Table("users").Select("id", "name")
	mainQB.With("cte_users", cteQB)

	cte := mainQB.GetData().WithCTEs[0]
	assert.Len(t, mainQB.GetData().WithCTEs, 1)
	assert.Equal(t, "cte_users", cte.Name)
	assert.NotNil(t, cte.Query)
	assert.Nil(t, cte.Expression)
	assert.False(t, cte.Recursive)

	sql, bindings, _ := mainQB.Select("*").ToSQL()
	expected := "WITH cte_users AS (SELECT id, name FROM users) SELECT *"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any(nil), bindings)
}

func Test_CTE_WithExpression(t *testing.T) {
	mainQB := xqb.Query()
	mainQB.WithExpr("cte_expr", "SELECT ?", 42)

	cte := mainQB.GetData().WithCTEs[0]
	assert.Len(t, mainQB.GetData().WithCTEs, 1)
	assert.Equal(t, "cte_expr", cte.Name)
	assert.Nil(t, cte.Query)
	assert.NotNil(t, cte.Expression)
	assert.Equal(t, "SELECT ?", cte.Expression.SQL)
	assert.Equal(t, []any{42}, cte.Expression.Bindings)

	sql, bindings, _ := mainQB.Select("*").ToSQL()
	expected := "WITH cte_expr AS (SELECT ?) SELECT *"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{42}, bindings)
}

func Test_CTE_WithRecursive(t *testing.T) {
	mainQB := xqb.Query()
	recQB := xqb.Table("tree").Select("id", "parent_id")
	mainQB.WithRecursive("cte_tree", recQB)

	cte := mainQB.GetData().WithCTEs[0]
	assert.True(t, cte.Recursive)

	sql, _, _ := mainQB.Select("*").ToSQL()
	expected := "WITH RECURSIVE cte_tree AS (SELECT id, parent_id FROM tree) SELECT *"
	assert.Equal(t, expected, sql)
}

func Test_CTE_WithRaw(t *testing.T) {
	mainQB := xqb.New()
	mainQB.WithRaw("cte_raw", "SELECT ? AS col", 99)

	cte := mainQB.GetData().WithCTEs[0]
	assert.NotNil(t, cte.Expression)
	assert.False(t, cte.Recursive)

	sql, bindings, _ := mainQB.Select("*").ToSQL()
	expected := "WITH cte_raw AS (SELECT ? AS col) SELECT *"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{99}, bindings)
}

func Test_CTE_WithRecursiveRaw(t *testing.T) {
	mainQB := xqb.Query()
	mainQB.WithRecursiveRaw("cte_rec_raw", "SELECT ? AS col", 123)

	cte := mainQB.GetData().WithCTEs[0]
	assert.True(t, cte.Recursive)

	sql, bindings, _ := mainQB.Select("*").ToSQL()
	expected := "WITH RECURSIVE cte_rec_raw AS (SELECT ? AS col) SELECT *"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{123}, bindings)
}

func Test_CTE_WithAdvancedExpressions(t *testing.T) {
	cteQB := xqb.Table("coverage_table").Select(
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

	mainQB := xqb.New()
	mainQB.With("cte_agg", cteQB).Select("*")

	sql, bindings, _ := mainQB.ToSQL()
	expected := "WITH cte_agg AS (SELECT SUM(amount) AS total_amount, LENGTH(bio) AS bio_len FROM coverage_table WHERE LOWER(status) = ? GROUP BY DATE(created_at), UPPER(region) HAVING total_amount > ? ORDER BY LENGTH(bio) DESC LIMIT 5 OFFSET 10) SELECT *"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{"active", 1000}, bindings)
}

func Test_CTE_WithMultipleCTEs(t *testing.T) {
	mainQB := xqb.Query()

	mainQB.
		WithRaw("cte1", "SELECT 1 AS one").
		WithExpr("cte2", "SELECT 2 AS two").
		With("cte3", xqb.Table("users").Select("id"))

	sql, _, _ := mainQB.Select("*").ToSQL()
	expected := "WITH cte1 AS (SELECT 1 AS one), cte2 AS (SELECT 2 AS two), cte3 AS (SELECT id FROM users) SELECT *"
	assert.Equal(t, expected, sql)
}

func Test_CTE_WithAliasedExpressions(t *testing.T) {
	mainQB := xqb.Query()
	mainQB.WithExpr("agg_stats", "SELECT COUNT(*) AS total, MAX(score) AS high_score FROM games")

	sql, _, _ := mainQB.Select("*").ToSQL()
	expected := "WITH agg_stats AS (SELECT COUNT(*) AS total, MAX(score) AS high_score FROM games) SELECT *"
	assert.Equal(t, expected, sql)
}

func Test_CTE_UsageInMainQuery(t *testing.T) {
	mainQB := xqb.New()
	mainQB.
		With("cte_users", xqb.New().Table("users").Select("id", "name")).
		From("cte_users").
		Where("id", ">", 5)

	sql, bindings, _ := mainQB.ToSQL()
	expected := "WITH cte_users AS (SELECT id, name FROM users) SELECT * FROM cte_users WHERE id > ?"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{5}, bindings)
}

func Test_CTE_Recursive_Usage(t *testing.T) {
	recQB := xqb.Table("tree").Select("id", "parent_id")
	mainQB := xqb.Table("tree_cte").
		WithRecursive("tree_cte", recQB).
		WhereNull("parent_id")

	sql, _, _ := mainQB.ToSQL()
	expected := "WITH RECURSIVE tree_cte AS (SELECT id, parent_id FROM tree) SELECT * FROM tree_cte WHERE parent_id IS NULL"
	assert.Equal(t, expected, sql)
}

func Test_CTE_BindingsOrder(t *testing.T) {
	mainQB := xqb.Query()

	mainQB.
		WithRaw("cte1", "SELECT ? AS one", 1).
		WithRaw("cte2", "SELECT ? AS two", 2).
		From("cte2").
		Where("two", ">", 3)

	sql, bindings, _ := mainQB.ToSQL()
	expected := "WITH cte1 AS (SELECT ? AS one), cte2 AS (SELECT ? AS two) SELECT * FROM cte2 WHERE two > ?"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{1, 2, 3}, bindings)
}

func Test_CTE_EmptyCTEsShouldNotEmitWith(t *testing.T) {
	qb := xqb.Table("users").Select("id")
	sql, bindings, _ := qb.ToSQL()
	expected := "SELECT id FROM users"
	assert.Equal(t, expected, sql)
	assert.Empty(t, bindings)
}

func Test_CTE_ComplexThreeLevelChain(t *testing.T) {
	mainQB := xqb.Query()

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

	sql, bindings, _ := mainQB.ToSQL()
	expected := "WITH " +
		"high_value_orders AS (SELECT user_id, total FROM orders WHERE total > ?), " +
		"user_order_details AS (SELECT high_value_orders.user_id, users.name FROM high_value_orders JOIN users ON users.id = high_value_orders.user_id), " +
		"user_order_summary AS (SELECT name, COUNT(*) AS order_count FROM user_order_details GROUP BY name) " +
		"SELECT * FROM user_order_summary WHERE order_count > ? ORDER BY order_count DESC"

	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{100, 5}, bindings)
}
