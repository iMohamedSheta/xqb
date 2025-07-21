package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/stretchr/testify/assert"
)

func Test_Select(t *testing.T) {
	qb := xqb.Table("users")
	qb.Select("id", "name", "email")
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT id, name, email FROM users", sql)
	assert.Empty(t, bindings)
}

func Test_Select_WithWhere(t *testing.T) {
	qb := xqb.Table("users")
	qb.Select("id", "name")
	qb.Where("age", ">", 18)
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT id, name FROM users WHERE age > ?", sql)
	assert.Equal(t, []any{18}, bindings)
}

func Test_Select_WithJoins(t *testing.T) {
	qb := xqb.Table("users")
	qb.Select("users.id", "users.name", "orders.id as order_id")
	qb.Join("orders", "users.id = orders.user_id")
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT users.id, users.name, orders.id as order_id FROM users JOIN orders ON users.id = orders.user_id", sql)
	assert.Empty(t, bindings)
}

func Test_Select_WithLeftJoins(t *testing.T) {
	qb := xqb.Table("users")
	qb.Select("users.id", "users.name", "orders.id as order_id").Where("users.id", ">", 55)
	qb.Join("orders", "users.id = orders.user_id").Where("orders.id", ">", 11)
	qb.LeftJoin("products", "orders.product_id = products.id")
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT users.id, users.name, orders.id as order_id FROM users JOIN orders ON users.id = orders.user_id LEFT JOIN products ON orders.product_id = products.id WHERE users.id > ? AND orders.id > ?", sql)
	assert.Equal(t, []any{55, 11}, bindings)
}

func Test_Select_WithGroupBy(t *testing.T) {
	qb := xqb.Table("orders")
	qb.Select("user_id", "COUNT(*) as order_count")
	qb.GroupBy("user_id")
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT user_id, COUNT(*) as order_count FROM orders GROUP BY user_id", sql)
	assert.Empty(t, bindings)
}

func Test_Select_WithHaving(t *testing.T) {
	qb := xqb.Table("orders")
	qb.Select("user_id", "COUNT(*) as order_count")
	qb.GroupBy("user_id")
	qb.Having("order_count", ">", 5)
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT user_id, COUNT(*) as order_count FROM orders GROUP BY user_id HAVING order_count > ?", sql)
	assert.Equal(t, []any{5}, bindings)
}

func Test_Select_WithOrderBy(t *testing.T) {
	qb := xqb.Table("users")
	qb.Select("id", "name")
	qb.OrderBy("name", "ASC")
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT id, name FROM users ORDER BY name ASC", sql)
	assert.Empty(t, bindings)
}

func Test_Select_WithLimitOffset(t *testing.T) {
	qb := xqb.Table("users")
	qb.Select("id", "name")
	qb.Limit(10)
	qb.Offset(20)
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT id, name FROM users LIMIT 10 OFFSET 20", sql)
	assert.Empty(t, bindings)
}

func Test_Select_WithAggregateFunctions(t *testing.T) {
	qb := xqb.Table("orders").
		Select(
			xqb.Sum("amount", "total_amount"),
			xqb.Avg("amount", "average_amount"),
			xqb.Count("id", "order_count"),
		)

	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT SUM(amount) AS total_amount, AVG(amount) AS average_amount, COUNT(id) AS order_count FROM orders", sql)
	assert.Empty(t, bindings)
}

func Test_Select_WithCTE(t *testing.T) {
	qb := xqb.Table("users")
	qb.WithRaw("user_totals", "SELECT user_id, SUM(amount) as total_spent FROM orders GROUP BY user_id")
	qb.Select("users.id", "users.name", "user_totals.total_spent")
	qb.Join("user_totals", "users.id = user_totals.user_id")
	sql, bindings, _ := qb.ToSQL()

	expectedSQL := "WITH user_totals AS (SELECT user_id, SUM(amount) as total_spent FROM orders GROUP BY user_id) " +
		"SELECT users.id, users.name, user_totals.total_spent FROM users JOIN user_totals ON users.id = user_totals.user_id"
	assert.Equal(t, expectedSQL, sql)
	assert.Empty(t, bindings)
}

func Test_Select_WithComplexCTE(t *testing.T) {
	qb := xqb.Table("products")
	qb.WithRaw("active_users",
		"WITH user_orders AS (SELECT user_id, COUNT(*) as order_count FROM orders GROUP BY user_id) "+
			"SELECT users.id, users.name, user_orders.order_count FROM users "+
			"JOIN user_orders ON users.id = user_orders.user_id")
	qb.Select("products.id", "products.name", "active_users.name as buyer")
	qb.Join("active_users", "products.id = active_users.id")
	sql, bindings, _ := qb.ToSQL()

	expectedSQL := "WITH active_users AS (WITH user_orders AS (SELECT user_id, COUNT(*) as order_count FROM orders GROUP BY user_id) " +
		"SELECT users.id, users.name, user_orders.order_count FROM users JOIN user_orders ON users.id = user_orders.user_id) " +
		"SELECT products.id, products.name, active_users.name as buyer FROM products JOIN active_users ON products.id = active_users.id"
	assert.Equal(t, expectedSQL, sql)
	assert.Empty(t, bindings)
}

func Test_Select_WithJSONExpressions(t *testing.T) {
	qb := xqb.Table("users")
	qb.Select(
		"id",
		"name",
		xqb.JsonExtract("metadata", "preferences.theme", "theme"),
	)
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT id, name, JSON_EXTRACT(metadata, '$.preferences.theme') AS theme FROM users", sql)
	assert.Empty(t, bindings)
}

func Test_Select_WithStringFunctions(t *testing.T) {
	qb := xqb.Table("users")
	qb.Select(
		"id",
		xqb.Concat([]string{
			"first_name",
			"' '",
			"last_name",
		}, "full_name"),
	)
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT id, CONCAT(first_name, ' ', last_name) AS full_name FROM users", sql)
	assert.Equal(t, []any(nil), bindings)
}

func Test_Select_WithDateFunctions(t *testing.T) {
	qb := xqb.Table("orders")
	qb.Select(
		"id",
		xqb.DateFormat("created_at", "%Y-%m-%d", "order_date"),
	)
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT id, DATE_FORMAT(created_at, '%Y-%m-%d') AS order_date FROM orders", sql)
	assert.Equal(t, []any(nil), bindings)
}

func Test_Select_WithMathExpressions(t *testing.T) {
	qb := xqb.Table("orders")
	qb.Select(
		"id",
		xqb.Math("amount * 1.1", "total_with_tax"),
	)
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT id, amount * 1.1 AS total_with_tax FROM orders", sql)
	assert.Empty(t, bindings)
}

func Test_Select_WithLocking(t *testing.T) {
	qb := xqb.Table("users")
	qb.Select("id", "name")
	qb.LockForUpdate()
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT id, name FROM users FOR UPDATE", sql)
	assert.Empty(t, bindings)
}

func Test_Select_WithUnion(t *testing.T) {
	qb := xqb.Table("users").Select("id", "name")
	qb.UnionRaw("SELECT id, name FROM users WHERE type = ?", "admin")
	qb.UnionRaw("SELECT id, name FROM users WHERE type = ?", "superuser")
	qb.UnionRaw("SELECT id, name FROM users WHERE type = ?", "guest")

	sql, bindings, _ := qb.ToSQL()

	expectedSQL := "(SELECT id, name FROM users) UNION (SELECT id, name FROM users WHERE type = ?) UNION (SELECT id, name FROM users WHERE type = ?) UNION (SELECT id, name FROM users WHERE type = ?)"
	expectedBindings := []any{"admin", "superuser", "guest"}

	assert.Equal(t, expectedSQL, sql)
	assert.Equal(t, expectedBindings, bindings)
}

func Test_Select_WithDistinct(t *testing.T) {
	qb := xqb.Table("users")
	qb.Select("name")
	qb.Distinct()
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT DISTINCT name FROM users", sql)
	assert.Empty(t, bindings)
}

func Test_Select_WithRawExpressions(t *testing.T) {
	qb := xqb.Table("users")
	sql, _, _ := qb.Select(
		xqb.Raw("COUNT(*) as total"),
		"name",
		xqb.Raw("CONCAT(first_name, ' ', last_name) as full_name"),
	).ToSQL()

	expected := "SELECT COUNT(*) as total, name, CONCAT(first_name, ' ', last_name) as full_name FROM users"
	assert.Equal(t, expected, sql)
}

func Test_Select_WithDateExpressions(t *testing.T) {
	qb := xqb.Table("orders")
	sql, _, _ := qb.Select(
		xqb.Raw("DATE_FORMAT(created_at, '%Y-%m') as month"),
		xqb.Raw("COUNT(*) as total_orders"),
		xqb.Raw("SUM(amount) as total_amount"),
	).
		GroupBy(xqb.Raw("DATE_FORMAT(created_at, '%Y-%m')")).
		OrderBy(xqb.Raw("DATE_FORMAT(created_at, '%Y-%m')"), "ASC").
		ToSQL()

	expected := "SELECT DATE_FORMAT(created_at, '%Y-%m') as month, COUNT(*) as total_orders, SUM(amount) as total_amount FROM orders GROUP BY DATE_FORMAT(created_at, '%Y-%m') ORDER BY DATE_FORMAT(created_at, '%Y-%m') ASC"
	assert.Equal(t, expected, sql)
}

func Test_Select_WithExpressions(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.Select(
		"id",
		xqb.Raw("CONCAT(first_name, ' ', last_name) as full_name"),
		xqb.Raw("(SELECT COUNT(*) FROM orders WHERE orders.user_id = users.id) as order_count"),
	).
		Where(xqb.Raw("LOWER(email)"), "LIKE", "%@example.com").
		GroupBy("id", "first_name", "last_name").
		Having(xqb.Raw("(SELECT COUNT(*) FROM orders WHERE orders.user_id = users.id)"), ">", 5).
		OrderBy(xqb.Raw("(SELECT SUM(amount) FROM orders WHERE orders.user_id = users.id)"), "DESC").
		ToSQL()

	expected := "SELECT id, CONCAT(first_name, ' ', last_name) as full_name, (SELECT COUNT(*) FROM orders WHERE orders.user_id = users.id) as order_count FROM users WHERE LOWER(email) LIKE ? GROUP BY id, first_name, last_name HAVING (SELECT COUNT(*) FROM orders WHERE orders.user_id = users.id) > ? ORDER BY (SELECT SUM(amount) FROM orders WHERE orders.user_id = users.id) DESC"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{"%@example.com", 5}, bindings)
}

func Test_Select_WithSubQuery(t *testing.T) {
	subSql, subBindings, _ := xqb.Table("payments").
		Select("id", "amount", "created_at").
		Where("payments.user_id", "=", 15).
		ToSQL()

	qb := xqb.Table("users").
		Select("id", "name", xqb.Raw("("+subSql+") AS payments", subBindings...)).
		Where("id", "=", 15)

	sql, bindings, _ := qb.ToSQL()

	expected := "SELECT id, name, (SELECT id, amount, created_at FROM payments WHERE payments.user_id = ?) AS payments FROM users WHERE id = ?"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{15, 15}, bindings)
}

func Test_Select_WithSubQuery_(t *testing.T) {
	sub := xqb.Table("payments").
		Select("id", "amount", "created_at").
		Where("payments.user_id", "=", 15)

	sub2 := xqb.Table("admins").
		Select("id", "amount", "created_at").
		Where("admins.user_id", "=", 15)

	qb := xqb.Table("users").
		Select("id", "name").
		SelectSub(sub, "payments").
		SelectSub(sub2, "admins").
		Where("id", "=", 15)

	sql, bindings, _ := qb.ToSQL()

	expected := "SELECT id, name, " +
		"(SELECT id, amount, created_at FROM payments WHERE payments.user_id = ?) AS payments, " +
		"(SELECT id, amount, created_at FROM admins WHERE admins.user_id = ?) AS admins " +
		"FROM users WHERE id = ?"

	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{15, 15, 15}, bindings)
}
func Test_FromSubquery(t *testing.T) {
	sub := xqb.Table("orders").
		Select("user_id", xqb.Raw("COUNT(*) AS order_count")).
		Where("user_id", "=", 25).
		GroupBy("user_id")

	qb := xqb.New().
		Select("u.id", "u.name", "o.order_count").
		FromSubquery(sub, "o").
		Join("users u", "u.id = o.user_id").
		Where("u.id", "=", 25)

	sql, bindings, _ := qb.ToSQL()

	expected := "SELECT u.id, u.name, o.order_count FROM (SELECT user_id, COUNT(*) AS order_count FROM orders WHERE user_id = ? GROUP BY user_id) AS o JOIN users u ON u.id = o.user_id WHERE u.id = ?"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{25, 25}, bindings)
}
