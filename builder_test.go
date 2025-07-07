package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/stretchr/testify/assert"
)

func TestBasicSelect(t *testing.T) {
	qb := xqb.Table("users")
	qb.Select("id", "name", "email")
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT id, name, email FROM users", sql)
	assert.Empty(t, bindings)
}

func TestSelectWithWhere(t *testing.T) {
	qb := xqb.Table("users")
	qb.Select("id", "name")
	qb.Where("age", ">", 18)
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT id, name FROM users WHERE age > ?", sql)
	assert.Equal(t, []interface{}{18}, bindings)
}

func TestSelectWithJoins(t *testing.T) {
	qb := xqb.Table("users")
	qb.Select("users.id", "users.name", "orders.id as order_id")
	qb.Join("orders", "users.id = orders.user_id")
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT users.id, users.name, orders.id as order_id FROM users INNER JOIN orders ON users.id = orders.user_id", sql)
	assert.Empty(t, bindings)
}

func TestSelectWithLeftJoins(t *testing.T) {
	qb := xqb.Table("users")
	qb.Select("users.id", "users.name", "orders.id as order_id").Where("users.id", ">", 55)
	qb.Join("orders", "users.id = orders.user_id").Where("orders.id", ">", 11)
	qb.LeftJoin("products", "orders.product_id = products.id")
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT users.id, users.name, orders.id as order_id FROM users INNER JOIN orders ON users.id = orders.user_id LEFT JOIN products ON orders.product_id = products.id WHERE users.id > ? AND orders.id > ?", sql)
	assert.Equal(t, []interface{}{55, 11}, bindings)
}

func TestSelectWithComplexJoins(t *testing.T) {

}

func TestSelectWithGroupBy(t *testing.T) {
	qb := xqb.Table("orders")
	qb.Select("user_id", "COUNT(*) as order_count")
	qb.GroupBy("user_id")
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT user_id, COUNT(*) as order_count FROM orders GROUP BY user_id", sql)
	assert.Empty(t, bindings)
}

func TestSelectWithHaving(t *testing.T) {
	qb := xqb.Table("orders")
	qb.Select("user_id", "COUNT(*) as order_count")
	qb.GroupBy("user_id")
	qb.Having("order_count", ">", 5)
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT user_id, COUNT(*) as order_count FROM orders GROUP BY user_id HAVING order_count > ?", sql)
	assert.Equal(t, []interface{}{5}, bindings)
}

func TestSelectWithOrderBy(t *testing.T) {
	qb := xqb.Table("users")
	qb.Select("id", "name")
	qb.OrderBy("name", "ASC")
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT id, name FROM users ORDER BY name ASC", sql)
	assert.Empty(t, bindings)
}

func TestSelectWithLimitOffset(t *testing.T) {
	qb := xqb.Table("users")
	qb.Select("id", "name")
	qb.Limit(10)
	qb.Offset(20)
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT id, name FROM users LIMIT 10 OFFSET 20", sql)
	assert.Empty(t, bindings)
}

func TestSelectWithAggregateFunctions(t *testing.T) {
	qb := xqb.Table("orders")
	qb.Sum("amount", "total_amount")
	qb.Avg("amount", "average_amount")
	qb.CountAggregate("id", "order_count")
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT SUM(amount) AS total_amount, AVG(amount) AS average_amount, COUNT(id) AS order_count FROM orders", sql)
	assert.Empty(t, bindings)
}

func TestSelectWithCTE(t *testing.T) {
	// Main query using the CTE
	qb := xqb.Table("users")
	qb.WithExpression("user_totals", "SELECT user_id, SUM(amount) as total_spent FROM orders GROUP BY user_id")
	qb.Select("users.id", "users.name", "user_totals.total_spent")
	qb.Join("user_totals", "users.id = user_totals.user_id")
	sql, bindings, _ := qb.ToSQL()

	expectedSQL := "WITH user_totals AS (SELECT user_id, SUM(amount) as total_spent FROM orders GROUP BY user_id) " +
		"SELECT users.id, users.name, user_totals.total_spent FROM users INNER JOIN user_totals ON users.id = user_totals.user_id"
	assert.Equal(t, expectedSQL, sql)
	assert.Empty(t, bindings)
}

func TestSelectWithComplexCTE(t *testing.T) {
	// Main query using nested CTEs
	qb := xqb.Table("products")
	qb.WithExpression("active_users",
		"WITH user_orders AS (SELECT user_id, COUNT(*) as order_count FROM orders GROUP BY user_id) "+
			"SELECT users.id, users.name, user_orders.order_count FROM users "+
			"INNER JOIN user_orders ON users.id = user_orders.user_id")
	qb.Select("products.id", "products.name", "active_users.name as buyer")
	qb.Join("active_users", "products.id = active_users.id")
	sql, bindings, _ := qb.ToSQL()

	expectedSQL := "WITH active_users AS (WITH user_orders AS (SELECT user_id, COUNT(*) as order_count FROM orders GROUP BY user_id) " +
		"SELECT users.id, users.name, user_orders.order_count FROM users INNER JOIN user_orders ON users.id = user_orders.user_id) " +
		"SELECT products.id, products.name, active_users.name as buyer FROM products INNER JOIN active_users ON products.id = active_users.id"
	assert.Equal(t, expectedSQL, sql)
	assert.Empty(t, bindings)
}

func TestSelectWithJSONExpressions(t *testing.T) {
	qb := xqb.Table("users")
	qb.Select("id", "name")
	qb.JSON("metadata", "$.preferences.theme", "theme")
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT id, name, JSON_EXTRACT(metadata, '$.preferences.theme') AS theme FROM users", sql)
	assert.Empty(t, bindings)
}

func TestSelectWithStringFunctions(t *testing.T) {
	qb := xqb.Table("users")
	qb.Select("id")
	qb.String("CONCAT", "first_name", []interface{}{" ", "last_name"}, "full_name")
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT id, CONCAT(first_name, ?, ?) AS full_name FROM users", sql)
	assert.Equal(t, []interface{}{" ", "last_name"}, bindings)
}

func TestSelectWithDateFunctions(t *testing.T) {
	qb := xqb.Table("orders")
	qb.Select("id")
	qb.Date("DATE_FORMAT", "created_at", []interface{}{"%Y-%m-%d"}, "order_date")
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT id, DATE_FORMAT(created_at, ?) AS order_date FROM orders", sql)
	assert.Equal(t, []interface{}{"%Y-%m-%d"}, bindings)
}

func TestSelectWithMathExpressions(t *testing.T) {
	qb := xqb.Table("orders")
	qb.Select("id")
	qb.Math("amount * 1.1", "total_with_tax")
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT id, amount * 1.1 AS total_with_tax FROM orders", sql)
	assert.Empty(t, bindings)
}

func TestSelectWithConditionalExpressions(t *testing.T) {
	qb := xqb.Table("orders")
	qb.Select("id")
	qb.Conditional("CASE WHEN status = 'completed' THEN 'done' ELSE 'pending' END", "status_text")
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT id, CASE WHEN status = 'completed' THEN 'done' ELSE 'pending' END AS status_text FROM orders", sql)
	assert.Empty(t, bindings)
}

func TestSelectWithIndexHints(t *testing.T) {
	qb := xqb.Table("users")
	qb.Select("id", "name")
	qb.ForceIndex("idx_name")
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT id, name FROM users FORCE INDEX (idx_name)", sql)
	assert.Empty(t, bindings)
}

func TestSelectWithLocking(t *testing.T) {
	qb := xqb.Table("users")
	qb.Select("id", "name")
	qb.ForUpdate()
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT id, name FROM users FOR UPDATE", sql)
	assert.Empty(t, bindings)
}

func TestSelectWithUnion(t *testing.T) {
	qb := xqb.Table("users").Select("id", "name")
	qb.Union("SELECT id, name FROM users WHERE type = ?", "admin")
	qb.Union("SELECT id, name FROM users WHERE type = ?", "superuser")
	qb.Union("SELECT id, name FROM users WHERE type = ?", "guest")

	sql, bindings, _ := qb.ToSQL()

	expectedSQL := "SELECT id, name FROM users UNION (SELECT id, name FROM users WHERE type = ?) UNION (SELECT id, name FROM users WHERE type = ?) UNION (SELECT id, name FROM users WHERE type = ?)"
	expectedBindings := []interface{}{"admin", "superuser", "guest"}

	assert.Equal(t, expectedSQL, sql)
	assert.Equal(t, expectedBindings, bindings)
}

func TestSelectWithDistinct(t *testing.T) {
	qb := xqb.Table("users")
	qb.Select("name")
	qb.Distinct()
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT DISTINCT name FROM users", sql)
	assert.Empty(t, bindings)
}

func TestSelectWithHighPriority(t *testing.T) {
	qb := xqb.Table("users")
	qb.Select("id", "name")
	qb.HighPriority()
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT HIGH_PRIORITY id, name FROM users", sql)
	assert.Empty(t, bindings)
}

func TestSelectWithStraightJoin(t *testing.T) {
	qb := xqb.Table("users")
	qb.Select("id", "name")
	qb.StraightJoin()
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT STRAIGHT_JOIN id, name FROM users", sql)
	assert.Empty(t, bindings)
}

func TestSelectWithCalcFoundRows(t *testing.T) {
	qb := xqb.Table("users")
	qb.Select("id", "name")
	qb.CalcFoundRows()
	sql, bindings, _ := qb.ToSQL()

	assert.Equal(t, "SELECT SQL_CALC_FOUND_ROWS id, name FROM users", sql)
	assert.Empty(t, bindings)
}

func TestSelectWithRawExpressions(t *testing.T) {
	qb := xqb.Table("users")
	sql, _, _ := qb.Select(xqb.Raw("COUNT(*) as total"), "name", xqb.Raw("CONCAT(first_name, ' ', last_name) as full_name")).ToSQL()

	expected := "SELECT COUNT(*) as total, name, CONCAT(first_name, ' ', last_name) as full_name FROM users"
	assert.Equal(t, expected, sql)
}

func TestWhereWithRawExpressions(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.Where(xqb.Raw("LOWER(name)"), "=", "john").ToSQL()

	expected := "SELECT * FROM users WHERE LOWER(name) = ?"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []interface{}{"john"}, bindings)
}

func TestWhereRaw(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.WhereRaw("LOWER(name) = ? OR LOWER(email) = ?", "john", "john@example.com").ToSQL()

	expected := "SELECT * FROM users WHERE LOWER(name) = ? OR LOWER(email) = ?"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []interface{}{"john", "john@example.com"}, bindings)
}

func TestGroupByWithRawExpressions(t *testing.T) {
	qb := xqb.Table("orders")
	sql, _, _ := qb.Select(xqb.Raw("YEAR(created_at) as year"), xqb.Raw("SUM(amount) as total")).
		GroupBy(xqb.Raw("YEAR(created_at)")).
		ToSQL()

	expected := "SELECT YEAR(created_at) as year, SUM(amount) as total FROM orders GROUP BY YEAR(created_at)"
	assert.Equal(t, expected, sql)
}

func TestOrderByWithRawExpressions(t *testing.T) {
	qb := xqb.Table("users")
	sql, _, _ := qb.OrderBy(xqb.Raw("FIELD(status, 'active', 'pending', 'inactive')"), "ASC").ToSQL()

	expected := "SELECT * FROM users ORDER BY FIELD(status, 'active', 'pending', 'inactive') ASC"
	assert.Equal(t, expected, sql)
}

func TestHavingWithRawExpressions(t *testing.T) {
	qb := xqb.Table("orders")
	sql, bindings, _ := qb.Select("user_id", xqb.Raw("SUM(amount) as total")).
		GroupBy("user_id").
		Having(xqb.Raw("SUM(amount)"), ">", 1000).
		ToSQL()

	expected := "SELECT user_id, SUM(amount) as total FROM orders GROUP BY user_id HAVING SUM(amount) > ?"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []interface{}{1000}, bindings)
}

func TestWhereNull(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.WhereNull("deleted_at").ToSQL()

	expected := "SELECT * FROM users WHERE deleted_at IS NULL"
	assert.Equal(t, expected, sql)
	assert.Empty(t, bindings)
}

func TestWhereNotNull(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.WhereNotNull("email").ToSQL()

	expected := "SELECT * FROM users WHERE email IS NOT NULL"
	assert.Equal(t, expected, sql)
	assert.Empty(t, bindings)
}

func TestComplexQueryWithExpressions(t *testing.T) {
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
	assert.Equal(t, []interface{}{"%@example.com", 5}, bindings)
}

func TestWhereWithSubquery(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.WhereRaw("EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id AND amount > ?)", 1000).ToSQL()

	expected := "SELECT * FROM users WHERE EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id AND amount > ?)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []interface{}{1000}, bindings)
}

func TestWhereWithCaseExpression(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.Where(xqb.Raw("CASE WHEN status = 'active' THEN 1 ELSE 0 END"), "=", 1).ToSQL()

	expected := "SELECT * FROM users WHERE CASE WHEN status = 'active' THEN 1 ELSE 0 END = ?"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []interface{}{1}, bindings)
}

func TestSelectWithDateExpressions(t *testing.T) {
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

func TestWhereGroup(t *testing.T) {
	qb := xqb.Table("orders")
	sql, _, _ := qb.WhereGroup(func(qb *xqb.QueryBuilder) {
		qb.Where("email", "=", "mohamed@mail.com").
			OrWhere("username", "=", "mohamed")
	}).WhereGroup(func(qb *xqb.QueryBuilder) {
		qb.Where("uuid", "=", "bbee7431-454d-4a8a-9435-961d191de2a7").OrWhere("user_id", "=", 4)
	}).OrWhereGroup(func(qb *xqb.QueryBuilder) {
		qb.Where("username", "=", "ahmed").Where("user_id", "=", 6)
	}).ToSQL()

	expected := "SELECT * FROM orders WHERE (email = ? OR username = ?) AND (uuid = ? OR user_id = ?) OR (username = ? AND user_id = ?)"
	assert.Equal(t, expected, sql)

}
