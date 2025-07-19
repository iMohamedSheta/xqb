package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/stretchr/testify/assert"
)

// .Where - Test cases for the Where clause in QueryBuilder

func Test_Where_Subquery_1(t *testing.T) {
	qb := xqb.Table("users")
	subQuery := xqb.Table("orders").Select("user_id").Where("status", "=", "active")
	sql, bindings, _ := qb.Where("id", "IN", subQuery).ToSQL()

	expected := "SELECT * FROM users WHERE id IN (SELECT user_id FROM orders WHERE status = ?)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{"active"}, bindings)
}

func Test_Where_Subquery_2(t *testing.T) {
	qb := xqb.Table("users")
	subQuery := xqb.Table("admins").Select("user_id").Where("role", "=", "superadmin").AddSelect("id").Latest("id")
	sql, bindings, _ := qb.Where("id", "IN", subQuery).ToSQL()
	expected := "SELECT * FROM users WHERE id IN (SELECT user_id, id FROM admins WHERE role = ? ORDER BY id DESC)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{"superadmin"}, bindings)
}

func Test_Where_Subquery_3(t *testing.T) {
	qb := xqb.Table("admins")
	subQuery := xqb.Table("users").
		Join("orders", "users.id = orders.user_id").
		Where("orders.status", "=", "paid").
		Select("users.id", "users.name", "orders.id AS order_id")

	sql, bindings, _ := qb.Where("id", "IN", subQuery).Where("admins.status", "=", "active").ToSQL()
	expected := "SELECT * FROM admins WHERE id IN (SELECT users.id, users.name, orders.id AS order_id FROM users JOIN orders ON users.id = orders.user_id WHERE orders.status = ?) AND admins.status = ?"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{"paid", "active"}, bindings)
}

func Test_Where_WithRaw_CaseExpression(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.Where(xqb.Raw("CASE WHEN status = 'active' THEN 1 ELSE 0 END"), "=", 1).ToSQL()

	expected := "SELECT * FROM users WHERE CASE WHEN status = 'active' THEN 1 ELSE 0 END = ?"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{1}, bindings)
}

func Test_Where_WithRaw_1(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.Join("orders", "users.id = orders.user_id").Where(xqb.Raw("CASE WHEN status = 'active' THEN 1 ELSE 0 END"), "=", 1).ToSQL()

	expected := "SELECT * FROM users JOIN orders ON users.id = orders.user_id WHERE CASE WHEN status = 'active' THEN 1 ELSE 0 END = ?"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{1}, bindings)
}

func Test_OrWhere_SubQuery_1(t *testing.T) {
	qb := xqb.Table("users")
	subQuery := xqb.Table("orders").Select("user_id").Where("status", "=", "active").Latest("id")
	sql, bindings, _ := qb.OrWhere("id", "IN", subQuery).ToSQL()
	expected := "SELECT * FROM users WHERE id IN (SELECT user_id FROM orders WHERE status = ? ORDER BY id DESC)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{"active"}, bindings)
}

func Test_OrWhere_SubQuery_2(t *testing.T) {
	qb := xqb.Table("users")
	subQuery := xqb.Table("orders").Join("admins", "users.id = admins.user_id").Select("user_id").Where("role", "=", "superadmin").Latest("id")
	sql, bindings, _ := qb.OrWhere("id", "IN", subQuery).ToSQL()
	expected := "SELECT * FROM users WHERE id IN (SELECT user_id FROM orders JOIN admins ON users.id = admins.user_id WHERE role = ? ORDER BY id DESC)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{"superadmin"}, bindings)
}

func Test_OrWhere_Raw_1(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.OrWhere(xqb.Raw("CASE WHEN status IN ('active', 'pending') THEN 1 ELSE 0 END"), "=", 1).ToSQL()
	expected := "SELECT * FROM users WHERE CASE WHEN status IN ('active', 'pending') THEN 1 ELSE 0 END = ?"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{1}, bindings)
}

func Test_OrWhere_Raw_2(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.OrWhere(xqb.Raw("CASE WHEN status IN ('active', 'pending') THEN 1 ELSE 0 END"), "=", 1).Join("orders", "users.id = orders.user_id").ToSQL()
	expected := "SELECT * FROM users JOIN orders ON users.id = orders.user_id WHERE CASE WHEN status IN ('active', 'pending') THEN 1 ELSE 0 END = ?"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{1}, bindings)
}

func Test_WhereNull_With_OrWhereNotNull(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.Where("id", "=", 1).WhereNull("deleted_at").OrWhereNull("disabled_at").ToSQL()
	expected := "SELECT * FROM users WHERE id = ? AND deleted_at IS NULL OR disabled_at IS NULL"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{1}, bindings)
}

func Test_WhereNull_With_Grouping(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.Where("id", "=", 1).WhereGroup(func(qb *xqb.QueryBuilder) {
		qb.WhereNull("deleted_at").OrWhereNull("disabled_at")
	}).ToSQL()

	expected := "SELECT * FROM users WHERE id = ? AND (deleted_at IS NULL OR disabled_at IS NULL)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, 1, len(bindings))
	assert.Equal(t, []any{1}, bindings)
}

func Test_WhereIn_normal(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.WhereIn("id", []any{1, 2, 3}).ToSQL()
	expected := "SELECT * FROM users WHERE id IN (?, ?, ?)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{1, 2, 3}, bindings)
}

func Test_WhereIn_With_Raw(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.WhereIn("id", []any{xqb.Raw("? UNION ?", 1, 2)}).ToSQL()
	expected := "SELECT * FROM users WHERE id IN (? UNION ?)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{1, 2}, bindings)
}

func Test_WhereIn_With_Raw_2(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.WhereIn("id", []any{xqb.Raw("? UNION ?", 1, 2)}).ToSQL()
	expected := "SELECT * FROM users WHERE id IN (? UNION ?)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{1, 2}, bindings)
}

func Test_WhereIn_With_Query(t *testing.T) {
	qb := xqb.Table("customers")
	subQuery := xqb.Table("users").Select("id").Where("type", "=", "active")
	sql, bindings, _ := qb.WhereIn("user_id", []any{subQuery}).ToSQL()
	expected := "SELECT * FROM customers WHERE user_id IN (SELECT id FROM users WHERE type = ?)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, 1, len(bindings))
	assert.Equal(t, []any{"active"}, bindings)
}

func Test_WhereIn_With_Query_Assert_If_There_Is_SubQuery_Use_It_Only(t *testing.T) {
	qb := xqb.Table("customers")
	subQuery := xqb.Table("users").Select("id").Where("type", "=", "active")
	sql, bindings, _ := qb.WhereIn("user_id", []any{15, 20, subQuery}).ToSQL()
	expected := "SELECT * FROM customers WHERE user_id IN (SELECT id FROM users WHERE type = ?)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, 1, len(bindings))
	assert.Equal(t, []any{"active"}, bindings)
}

func Test_WhereInQuery(t *testing.T) {
	qb := xqb.Table("customers")
	subQuery := xqb.Table("users").Select("id").Where("type", "=", "active")
	sql, bindings, _ := qb.WhereInQuery("user_id", subQuery).ToSQL()
	expected := "SELECT * FROM customers WHERE user_id IN (SELECT id FROM users WHERE type = ?)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, 1, len(bindings))
	assert.Equal(t, []any{"active"}, bindings)
}

func Test_WhereExists_With_SubQuery_1(t *testing.T) {
	qb := xqb.Table("users")
	subQuery := xqb.Table("admins").Select("user_id").Where("role", "IN", []any{"superadmin", "admin"}).Latest("id")
	sql, bindings, _ := qb.Select("1").WhereExists(subQuery).ToSQL()
	expected := "SELECT 1 FROM users WHERE EXISTS (SELECT user_id FROM admins WHERE role IN (?, ?) ORDER BY id DESC)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, 2, len(bindings))
}

func Test_WhereExists_With_SubQuery_2(t *testing.T) {
	qb := xqb.Table("customers")
	subQuery := xqb.Table("users").Select("id").Where("type", "=", "active")
	sql, bindings, _ := qb.Select("1").WhereExists(subQuery).ToSQL()
	expected := "SELECT 1 FROM customers WHERE EXISTS (SELECT id FROM users WHERE type = ?)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, 1, len(bindings))
}

func Test_WhereExists_With_Raw(t *testing.T) {
	qb := xqb.Table("orders")
	raw := xqb.Raw("SELECT user_id FROM users WHERE type = ?", "active")
	sql, bindings, _ := qb.Select("1").WhereExists(raw).ToSQL()
	expected := "SELECT 1 FROM orders WHERE EXISTS (SELECT user_id FROM users WHERE type = ?)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, 1, len(bindings))
	assert.Equal(t, []any{"active"}, bindings)
}

func Test_WhereNotExists_With_SubQuery_1(t *testing.T) {
	qb := xqb.Table("users")
	subQuery := xqb.Table("admins").Select("user_id").Where("role", "IN", []any{"superadmin", "admin"}).Latest("id")
	sql, bindings, _ := qb.Select("1").WhereNotExists(subQuery).ToSQL()
	expected := "SELECT 1 FROM users WHERE NOT EXISTS (SELECT user_id FROM admins WHERE role IN (?, ?) ORDER BY id DESC)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, 2, len(bindings))
}

func Test_OrWhereExists_WithSubQuery(t *testing.T) {
	qb := xqb.Table("users")
	subQuery := xqb.Table("admins").Select("user_id").Where("role", "IN", []any{"superadmin", "admin"}).Latest("id")
	sql, bindings, _ := qb.Select("1").Where("id", "=", 15).OrWhereExists(subQuery).ToSQL()
	expected := "SELECT 1 FROM users WHERE id = ? OR EXISTS (SELECT user_id FROM admins WHERE role IN (?, ?) ORDER BY id DESC)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, 3, len(bindings))
	assert.Equal(t, 15, bindings[0])
	assert.Equal(t, "superadmin", bindings[1])
	assert.Equal(t, "admin", bindings[2])
}

func Test_WhereValue(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.WhereValue("age", ">", 18).ToSQL()
	assert.Equal(t, "SELECT * FROM users WHERE age > ?", sql)
	assert.Equal(t, []any{18}, bindings)
}

func Test_OrWhereValue(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.Where("name", "=", "admin").OrWhereValue("role", "=", "guest").ToSQL()
	assert.Equal(t, "SELECT * FROM users WHERE name = ? OR role = ?", sql)
	assert.Equal(t, []any{"admin", "guest"}, bindings)
}

func Test_WhereExpr(t *testing.T) {
	expr := xqb.Raw("LOWER(name)")
	qb := xqb.Table("users")
	sql, bindings, _ := qb.WhereExpr("LOWER(name)", "=", expr).ToSQL()
	assert.Equal(t, "SELECT * FROM users WHERE LOWER(name) = (LOWER(name))", sql)
	assert.Len(t, bindings, 0)
}

func Test_OrWhereExpr(t *testing.T) {
	expr := xqb.Raw("LOWER(role)")
	qb := xqb.Table("users")
	sql, bindings, _ := qb.Where("name", "=", "mohamed").OrWhereExpr("LOWER(role)", "=", expr).ToSQL()
	assert.Equal(t, "SELECT * FROM users WHERE name = ? OR LOWER(role) = (LOWER(role))", sql)
	assert.Equal(t, []any{"mohamed"}, bindings)
}

func Test_WhereSub(t *testing.T) {
	sub := xqb.Table("admins").Select("id").Where("active", "=", true)
	qb := xqb.Table("users").WhereSub("admin_id", "IN", sub)
	sql, bindings, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM users WHERE admin_id IN (SELECT id FROM admins WHERE active = ?)", sql)
	assert.Equal(t, []any{true}, bindings)
}

func Test_OrWhereSub(t *testing.T) {
	sub := xqb.Table("admins").Select("id").Where("active", "=", true)
	qb := xqb.Table("users").Where("role", "=", "staff").OrWhereSub("admin_id", "IN", sub)
	sql, bindings, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM users WHERE role = ? OR admin_id IN (SELECT id FROM admins WHERE active = ?)", sql)
	assert.Equal(t, []any{"staff", true}, bindings)
}

func Test_WhereNotInQuery(t *testing.T) {
	sub := xqb.Table("banned_users").Select("id")
	qb := xqb.Table("users").WhereNotInQuery("id", sub)
	sql, _, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM users WHERE id NOT IN (SELECT id FROM banned_users)", sql)
}

func Test_OrWhereNotInQuery(t *testing.T) {
	sub := xqb.Table("banned_users").Select("id")
	qb := xqb.Table("users").Where("role", "=", "staff").OrWhereNotInQuery("id", sub)
	sql, bindings, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM users WHERE role = ? OR id NOT IN (SELECT id FROM banned_users)", sql)
	assert.Equal(t, []any{"staff"}, bindings)
}

func Test_WhereNotBetween(t *testing.T) {
	qb := xqb.Table("users").WhereNotBetween("age", 18, 60)
	sql, bindings, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM users WHERE age NOT BETWEEN ? AND ?", sql)
	assert.Equal(t, []any{18, 60}, bindings)
}

func Test_OrWhereNotBetween(t *testing.T) {
	qb := xqb.Table("users").Where("role", "=", "guest").OrWhereNotBetween("age", 10, 20)
	sql, bindings, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM users WHERE role = ? OR age NOT BETWEEN ? AND ?", sql)
	assert.Equal(t, []any{"guest", 10, 20}, bindings)
}

func Test_WhereGroup_MultipleLevels(t *testing.T) {
	qb := xqb.Table("users").WhereGroup(func(q1 *xqb.QueryBuilder) {
		q1.Where("status", "=", "active").OrWhereGroup(func(q2 *xqb.QueryBuilder) {
			q2.Where("email_verified", "=", false).Where("banned", "=", false)
		})
	})

	sql, bindings, _ := qb.ToSQL()
	expected := "SELECT * FROM users WHERE (status = ? OR (email_verified = ? AND banned = ?))"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{"active", false, false}, bindings)
}

func Test_WhereRaw_WithBindings(t *testing.T) {
	qb := xqb.Table("logs").WhereRaw("created_at > ?", "2024-01-01")
	sql, bindings, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM logs WHERE created_at > ?", sql)
	assert.Equal(t, []any{"2024-01-01"}, bindings)
}

func Test_OrWhereRaw_WithBindings(t *testing.T) {
	qb := xqb.Table("logs").
		Where("type", "=", "info").
		OrWhereRaw("created_at > ?", "2024-01-01")
	sql, bindings, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM logs WHERE type = ? OR created_at > ?", sql)
	assert.Equal(t, []any{"info", "2024-01-01"}, bindings)
}

func Test_WhereIn_Empty(t *testing.T) {
	qb := xqb.Table("users")
	sql, _, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM users", sql) // Nothing should be added
}

func Test_WhereNotIn_Empty(t *testing.T) {
	qb := xqb.Table("users").WhereNotIn("id", []any{})
	sql, _, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM users", sql) // Nothing should be added
}

func Test_WhereBetween_WithExpr(t *testing.T) {
	min := xqb.Raw("NOW() - INTERVAL 1 DAY")
	max := xqb.Raw("NOW()")
	qb := xqb.Table("logs").WhereBetween("created_at", min, max)
	sql, bindings, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM logs WHERE created_at BETWEEN NOW() - INTERVAL 1 DAY AND NOW()", sql)
	assert.Len(t, bindings, 0)
}

func Test_WhereExists_Chained(t *testing.T) {
	sub := xqb.Table("admins").Select("id").Where("active", "=", true)
	qb := xqb.Table("users").Where("status", "=", "staff").WhereExists(sub)
	sql, bindings, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM users WHERE status = ? AND EXISTS (SELECT id FROM admins WHERE active = ?)", sql)
	assert.Equal(t, []any{"staff", true}, bindings)
}

func Test_WhereNotExists_Chained(t *testing.T) {
	sub := xqb.Table("admins").Select("id").Where("active", "=", false)
	qb := xqb.Table("users").Where("status", "=", "guest").WhereNotExists(sub)
	sql, bindings, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM users WHERE status = ? AND NOT EXISTS (SELECT id FROM admins WHERE active = ?)", sql)
	assert.Equal(t, []any{"guest", false}, bindings)
}

func Test_Mixed_WhereRaw_And_Normal(t *testing.T) {
	qb := xqb.Table("users").
		WhereRaw("JSON_EXTRACT(meta, '$.age') > ?", 18).
		Where("active", "=", true)
	sql, bindings, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM users WHERE JSON_EXTRACT(meta, '$.age') > ? AND active = ?", sql)
	assert.Equal(t, []any{18, true}, bindings)
}

func Test_OrWhereGroup_Complex(t *testing.T) {
	qb := xqb.Table("products").
		Where("stock", ">", 0).
		OrWhereGroup(func(q *xqb.QueryBuilder) {
			q.Where("archived", "=", false).
				Where("discontinued", "=", false)
		})
	sql, bindings, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM products WHERE stock > ? OR (archived = ? AND discontinued = ?)", sql)
	assert.Equal(t, []any{0, false, false}, bindings)
}

func Test_WhereExpr_ComplexBothSides(t *testing.T) {
	left := xqb.Raw("LOWER(username)")
	right := xqb.Raw("LOWER(?)", "Mohamed")
	qb := xqb.Table("users").Where(left, "=", right)
	sql, bindings, _ := qb.ToSQL()
	assert.Equal(t, "SELECT * FROM users WHERE (LOWER(username)) = (LOWER(?))", sql)
	assert.Equal(t, []any{"Mohamed"}, bindings)
}

func TestWhereWithRawExpressions(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.Where(xqb.Raw("LOWER(name)"), "=", "john").ToSQL()

	expected := "SELECT * FROM users WHERE LOWER(name) = ?"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{"john"}, bindings)
}

func TestWhereRaw(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.WhereRaw("LOWER(name) = ? OR LOWER(email) = ?", "john", "john@example.com").ToSQL()

	expected := "SELECT * FROM users WHERE LOWER(name) = ? OR LOWER(email) = ?"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{"john", "john@example.com"}, bindings)
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

func TestWhereNullWithSelect(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.Select("id", "name").Where("name", "LIKE", "%mohamedsheta%").WhereNull("deleted_at").ToSQL()

	expected := "SELECT id, name FROM users WHERE name LIKE ? AND deleted_at IS NULL"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{"%mohamedsheta%"}, bindings)
}

func TestWhereNotNullWithSelect(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.Select("id", "name").Where("name", "LIKE", "%mohamedsheta%").WhereNotNull("email").ToSQL()

	expected := "SELECT id, name FROM users WHERE name LIKE ? AND email IS NOT NULL"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{"%mohamedsheta%"}, bindings)
}

func TestWhereIn(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.WhereIn("id", []any{1, 2, 3}).ToSQL()

	expected := "SELECT * FROM users WHERE id IN (?, ?, ?)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{1, 2, 3}, bindings)
}

func TestWhereNotIn(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.WhereNotIn("id", []any{1, 2, 3}).ToSQL()

	expected := "SELECT * FROM users WHERE id NOT IN (?, ?, ?)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{1, 2, 3}, bindings)
}

func TestWhereInWithSubquery(t *testing.T) {
	qb := xqb.Table("users")
	subQuery := xqb.Table("orders").Select("user_id").Where("status", "=", "active")
	sql, bindings, _ := qb.WhereIn("id", []any{subQuery}).ToSQL()

	expected := "SELECT * FROM users WHERE id IN (SELECT user_id FROM orders WHERE status = ?)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{"active"}, bindings)
}

func TestWhereBetween(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.WhereBetween("age", 18, 30).ToSQL()

	expected := "SELECT * FROM users WHERE age BETWEEN ? AND ?"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{18, 30}, bindings)
}

func TestWhereRawWithSubqueryRaw(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.WhereRaw("EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id AND amount > ?)", 1000).ToSQL()

	expected := "SELECT * FROM users WHERE EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id AND amount > ?)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{1000}, bindings)
}

func Test_WhereGroup(t *testing.T) {
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

func Test_Where_Is_Null(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.Where("id", "=", 1).Where("deleted_at", "IS NULL", nil).ToSQL()

	expected := "SELECT * FROM users WHERE id = ? AND deleted_at IS NULL"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{1}, bindings)
}
