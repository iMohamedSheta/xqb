package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/stretchr/testify/assert"
)

func Test_Join_String_Table(t *testing.T) {
	qb := xqb.Table("users").Join("posts", "users.id = posts.user_id")
	sql, bindings, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM users JOIN posts ON users.id = posts.user_id", sql)
	assert.Empty(t, bindings)
}

func Test_Join_With_Bindings(t *testing.T) {
	qb := xqb.Table("users").Join("posts", "users.id = posts.user_id AND posts.status = ?", "active")
	sql, bindings, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM users JOIN posts ON users.id = posts.user_id AND posts.status = ?", sql)
	assert.Equal(t, []any{"active"}, bindings)
}

func Test_LeftJoin(t *testing.T) {
	qb := xqb.Table("users").LeftJoin("comments", "users.id = comments.user_id")
	sql, _, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM users LEFT JOIN comments ON users.id = comments.user_id", sql)
}

func Test_RightJoin(t *testing.T) {
	qb := xqb.Table("users").RightJoin("logins", "users.id = logins.user_id")
	sql, _, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM users RIGHT JOIN logins ON users.id = logins.user_id", sql)
}

func Test_FullJoin(t *testing.T) {
	qb := xqb.Table("users").FullJoin("sessions", "users.id = sessions.user_id")
	sql, _, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM users FULL JOIN sessions ON users.id = sessions.user_id", sql)
}

func Test_CrossJoin(t *testing.T) {
	qb := xqb.Table("users").CrossJoin("roles")
	sql, _, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM users CROSS JOIN roles", sql)
}

func Test_CrossJoinSub(t *testing.T) {
	sub := xqb.Table("plans").Where("expired", "=", false)
	qb := xqb.Table("users").CrossJoinSub(sub, "p")
	sql, bindings, _ := qb.ToSQL()
	assert.Contains(t, sql, "CROSS JOIN (SELECT * FROM plans WHERE expired = ?) AS p")
	assert.Equal(t, []any{false}, bindings)
}

func Test_CrossJoin_With_Expr(t *testing.T) {
	raw := xqb.Raw("(SELECT * FROM regions WHERE active = ?) AS r", true)
	qb := xqb.Table("users").CrossJoinExpr(raw)
	sql, bindings, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM users CROSS JOIN (SELECT * FROM regions WHERE active = ?) AS r", sql)
	assert.Equal(t, []any{true}, bindings)
}

func Test_Join_SubQuery_DefaultAlias(t *testing.T) {
	sub := xqb.Table("posts").Where("published", "=", true)
	qb := xqb.Table("users").JoinSub(sub, "sub", "users.id = sub.user_id")
	sql, bindings, _ := qb.ToSQL()
	assert.Contains(t, sql, "JOIN (SELECT * FROM posts WHERE published = ?) AS sub ON users.id = sub.user_id")
	assert.Equal(t, []any{true}, bindings)
}

func Test_Join_SubQuery_With_Alias(t *testing.T) {
	sub := xqb.Table("posts").Where("published", "=", true)
	qb := xqb.Table("users").JoinSub(sub, "p", "users.id = p.user_id")
	sql, bindings, _ := qb.ToSQL()
	assert.Contains(t, sql, "JOIN (SELECT * FROM posts WHERE published = ?) AS p ON users.id = p.user_id")
	assert.Equal(t, []any{true}, bindings)
}

func Test_LeftJoin_SubQuery(t *testing.T) {
	sub := xqb.Table("comments").Where("active", "=", true)
	qb := xqb.Table("users").LeftJoinSub(sub, "c", "users.id = c.user_id")
	sql, bindings, _ := qb.ToSQL()
	assert.Contains(t, sql, "LEFT JOIN (SELECT * FROM comments WHERE active = ?) AS c ON users.id = c.user_id")
	assert.Equal(t, []any{true}, bindings)
}

func Test_RightJoin_SubQuery(t *testing.T) {
	sub := xqb.Table("orders").Where("status", "=", "paid")
	qb := xqb.Table("users").RightJoinSub(sub, "o", "users.id = o.user_id")
	sql, bindings, _ := qb.ToSQL()
	assert.Contains(t, sql, "RIGHT JOIN (SELECT * FROM orders WHERE status = ?) AS o ON users.id = o.user_id")
	assert.Equal(t, []any{"paid"}, bindings)
}

func Test_Join_With_Condition_Expression(t *testing.T) {
	qb := xqb.Table("users").Join("posts", "users.id = posts.user_id AND posts.status = ?", "active")
	sql, bindings, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM users JOIN posts ON users.id = posts.user_id AND posts.status = ?", sql)
	assert.Equal(t, []any{"active"}, bindings)
}

func Test_Join_With_Expression_Table(t *testing.T) {
	table := xqb.Raw("(SELECT * FROM posts WHERE published = ?) AS p", true)
	qb := xqb.Table("users").JoinExpr(table, "users.id = p.user_id")
	sql, bindings, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM users JOIN (SELECT * FROM posts WHERE published = ?) AS p ON users.id = p.user_id", sql)
	assert.Equal(t, []any{true}, bindings)
}

func Test_FullJoinExpr(t *testing.T) {
	expr := xqb.Raw("(SELECT * FROM stats WHERE active = ?) AS s", true)
	qb := xqb.Table("users").FullJoinExpr(expr, "users.id = s.user_id")
	sql, bindings, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM users FULL JOIN (SELECT * FROM stats WHERE active = ?) AS s ON users.id = s.user_id", sql)
	assert.Equal(t, []any{true}, bindings)
}

func Test_JoinSub_FallbackAlias(t *testing.T) {
	sub := xqb.Table("posts").Where("published", "=", true)
	qb := xqb.Table("users").JoinSub(sub, "", "users.id = sub.user_id")
	sql, _, _ := qb.ToSQL()
	assert.Contains(t, sql, "JOIN (SELECT * FROM posts WHERE published = ?) AS sub ON users.id = sub.user_id")
}

func Test_JoinExpr_With_Expression_Condition(t *testing.T) {
	table := xqb.Raw("(SELECT * FROM payments WHERE confirmed = ?) AS p", true)
	cond := xqb.Raw("users.id = p.user_id AND p.status = ?", "success")
	qb := xqb.Table("users").JoinExpr(table, cond)
	sql, bindings, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM users JOIN (SELECT * FROM payments WHERE confirmed = ?) AS p ON users.id = p.user_id AND p.status = ?", sql)
	assert.Equal(t, []any{true, "success"}, bindings)
}

func Test_CrossJoinSub_FallbackAlias(t *testing.T) {
	sub := xqb.Table("plans").Where("expired", "=", false)
	qb := xqb.Table("users").CrossJoinSub(sub, "")
	sql, bindings, _ := qb.ToSQL()
	assert.Contains(t, sql, "CROSS JOIN (SELECT * FROM plans WHERE expired = ?) AS sub")
	assert.Equal(t, []any{false}, bindings)
}

func Test_Multiple_Joins_Mixed_Types(t *testing.T) {
	sub := xqb.Table("orders").Where("status", "=", "shipped")
	expr := xqb.Raw("(SELECT * FROM invoices WHERE paid = ?) AS inv", true)
	qb := xqb.Table("users").
		Join("addresses", "users.id = addresses.user_id AND addresses.city = ?", "Cairo").
		LeftJoinSub(sub, "o", "users.id = o.user_id").
		RightJoinExpr(expr, "users.id = inv.user_id AND inv.total > ?", 1000)

	sql, bindings, _ := qb.ToSQL()
	assert.Contains(t, sql, "JOIN addresses ON users.id = addresses.user_id AND addresses.city = ?")
	assert.Contains(t, sql, "LEFT JOIN (SELECT * FROM orders WHERE status = ?) AS o ON users.id = o.user_id")
	assert.Contains(t, sql, "RIGHT JOIN (SELECT * FROM invoices WHERE paid = ?) AS inv ON users.id = inv.user_id AND inv.total > ?")
	assert.Equal(t, []any{"Cairo", "shipped", true, 1000}, bindings)
}

func Test_Join_With_SubQuery_That_Has_Join(t *testing.T) {
	innerSub := xqb.Table("payments").Where("amount", ">", 500)
	sub := xqb.Table("orders").
		JoinSub(innerSub, "pay", "orders.payment_id = pay.id").
		Where("orders.status", "=", "completed")

	qb := xqb.Table("users").JoinSub(sub, "o", "users.id = o.user_id")

	sql, bindings, _ := qb.ToSQL()
	assert.Contains(t, sql, "JOIN (SELECT * FROM orders JOIN (SELECT * FROM payments WHERE amount > ?) AS pay ON orders.payment_id = pay.id WHERE orders.status = ?) AS o ON users.id = o.user_id")
	assert.Equal(t, []any{500, "completed"}, bindings)
}

func Test_CrossJoin_Combined_With_Other_Joins(t *testing.T) {
	qb := xqb.Table("users").
		Join("posts", "users.id = posts.user_id").
		CrossJoin("countries")

	sql, bindings, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM users JOIN posts ON users.id = posts.user_id CROSS JOIN countries", sql)
	assert.Empty(t, bindings)
}

func Test_Stores_With_Orders_SubQuery(t *testing.T) {
	subQuery := xqb.Table("orders").
		Select("store_id", xqb.Raw("COUNT(*) as total_orders")).
		WhereNull("cancelled_at").
		WhereNotNull("confirmed_at").
		Where("status", "!=", "failed").
		GroupBy("store_id")

	sql, bindings, _ := xqb.Table("stores").
		Select(
			"managers.fullname",
			"managers.email",
			"stores.id",
			"order_stats.total_orders",
		).
		AddSelectRaw("locations.city location_city").
		AddSelectRaw("locations.zip_code location_zip").
		AddSelectRaw("managers.id manager_id").
		LeftJoinSub(subQuery, "order_stats", "stores.id = order_stats.store_id").
		Join("managers", "stores.manager_id = managers.id").
		Join("locations", "stores.location_id = locations.id").
		OrderBy("stores.id", "ASC").
		Where("stores.region_id", "=", 22).
		Limit(5).
		ToSQL()

	assert.Equal(t,
		"SELECT managers.fullname, managers.email, stores.id, order_stats.total_orders, locations.city location_city, locations.zip_code location_zip, managers.id manager_id "+
			"FROM stores "+
			"LEFT JOIN (SELECT store_id, COUNT(*) as total_orders FROM orders WHERE cancelled_at IS NULL AND confirmed_at IS NOT NULL AND status != ? GROUP BY store_id) AS order_stats ON stores.id = order_stats.store_id "+
			"JOIN managers ON stores.manager_id = managers.id "+
			"JOIN locations ON stores.location_id = locations.id "+
			"WHERE stores.region_id = ? "+
			"ORDER BY stores.id ASC "+
			"LIMIT 5",
		sql,
	)
	assert.Equal(t, []any{"failed", 22}, bindings)
}
