package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/stretchr/testify/assert"
)

func TestLimit(t *testing.T) {
	qb := xqb.Table("users").Select("*").Limit(10)
	sql, _, _ := qb.ToSQL()

	assert.Equal(t, "SELECT * FROM users LIMIT 10", sql)
}

func TestOffset(t *testing.T) {
	qb := xqb.Table("users").Select("*").Offset(5)
	sql, _, _ := qb.ToSQL()

	assert.Equal(t, "SELECT * FROM users OFFSET 5", sql)
}

func TestSkipAlias(t *testing.T) {
	qb := xqb.Table("users").Select("*").Skip(7)
	sql, _, _ := qb.ToSQL()

	assert.Equal(t, "SELECT * FROM users OFFSET 7", sql)
}

func TestTakeAlias(t *testing.T) {
	qb := xqb.Table("users").Select("*").Take(25)
	sql, _, _ := qb.ToSQL()

	assert.Equal(t, "SELECT * FROM users LIMIT 25", sql)
}

func TestForPage(t *testing.T) {
	qb := xqb.Table("users").Select("*").ForPage(3, 15) // page 3 = skip 30, take 15
	sql, _, _ := qb.ToSQL()

	assert.Equal(t, "SELECT * FROM users LIMIT 15 OFFSET 30", sql)
}

func TestLimitOffsetWithWhere(t *testing.T) {
	qb := xqb.Table("products").
		Select("id", "name").
		Where("price", ">", 100).
		OrderBy("created_at", "desc").
		Limit(20).
		Offset(40)

	sql, bindings, _ := qb.ToSQL()

	expectedSQL := "SELECT id, name FROM products WHERE price > ? ORDER BY created_at desc LIMIT 20 OFFSET 40"
	assert.Equal(t, expectedSQL, sql)
	assert.Equal(t, []any{100}, bindings)
}

func TestForPageWithWhereAndOrder(t *testing.T) {
	qb := xqb.Table("orders").
		Select("id", "user_id").
		Where("status", "=", "pending").
		OrderBy("id", "asc").
		ForPage(5, 10) // OFFSET = 40, LIMIT = 10

	sql, bindings, _ := qb.ToSQL()

	expectedSQL := "SELECT id, user_id FROM orders WHERE status = ? ORDER BY id asc LIMIT 10 OFFSET 40"
	assert.Equal(t, expectedSQL, sql)
	assert.Equal(t, []any{"pending"}, bindings)
}

func TestPaginationWithJoins(t *testing.T) {
	qb := xqb.Table("users").
		Select("users.id", "profiles.bio").
		Join("profiles", "profiles.user_id = users.id").
		OrderBy("users.created_at", "desc").
		Limit(50).
		Offset(100)

	sql, _, _ := qb.ToSQL()

	expectedSQL := "SELECT users.id, profiles.bio FROM users JOIN profiles ON profiles.user_id = users.id ORDER BY users.created_at desc LIMIT 50 OFFSET 100"
	assert.Equal(t, expectedSQL, sql)
}

func TestForPageLargePageNumber(t *testing.T) {
	qb := xqb.Table("logs").
		Select("*").
		ForPage(999, 1000) // OFFSET = 998000, LIMIT = 1000

	sql, _, _ := qb.ToSQL()

	expectedSQL := "SELECT * FROM logs LIMIT 1000 OFFSET 998000"
	assert.Equal(t, expectedSQL, sql)
}

func TestForPageWithGroupByHaving(t *testing.T) {
	qb := xqb.Table("transactions").
		Select("user_id", xqb.Raw("SUM(amount) as total")).
		GroupBy("user_id").
		Having("SUM(amount)", ">", 1000).
		ForPage(2, 25)

	sql, bindings, _ := qb.ToSQL()

	expectedSQL := "SELECT user_id, SUM(amount) as total FROM transactions GROUP BY user_id HAVING SUM(amount) > ? LIMIT 25 OFFSET 25"
	assert.Equal(t, expectedSQL, sql)
	assert.Equal(t, []any{1000}, bindings)
}
