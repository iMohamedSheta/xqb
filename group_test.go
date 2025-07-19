package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/stretchr/testify/assert"
)

func TestGroupByWithRawExpressions(t *testing.T) {
	qb := xqb.Table("orders")
	sql, _, _ := qb.Select(xqb.Raw("YEAR(created_at) as year"), xqb.Raw("SUM(amount) as total")).
		GroupBy(xqb.Raw("YEAR(created_at)")).
		ToSQL()

	expected := "SELECT YEAR(created_at) as year, SUM(amount) as total FROM orders GROUP BY YEAR(created_at)"
	assert.Equal(t, expected, sql)
}

func TestGroupByMultipleColumns(t *testing.T) {
	qb := xqb.Table("orders")
	sql, _, _ := qb.Select("user_id", "product_id").
		GroupBy("user_id", "product_id").
		ToSQL()

	expected := "SELECT user_id, product_id FROM orders GROUP BY user_id, product_id"
	assert.Equal(t, expected, sql)
}

func TestGroupByRawShortcut(t *testing.T) {
	qb := xqb.Table("orders")
	sql, _, _ := qb.Select("id").
		GroupByRaw("DATE(created_at)").
		ToSQL()

	expected := "SELECT id FROM orders GROUP BY DATE(created_at)"
	assert.Equal(t, expected, sql)
}

func TestGroupByWithHaving(t *testing.T) {
	qb := xqb.Table("sales")
	sql, bindings, _ := qb.
		Select("region", xqb.Raw("SUM(amount) as total")).
		GroupBy("region").
		Having("SUM(amount)", ">", 1000).
		ToSQL()

	expectedSQL := "SELECT region, SUM(amount) as total FROM sales GROUP BY region HAVING SUM(amount) > ?"
	assert.Equal(t, expectedSQL, sql)
	assert.Equal(t, []any{1000}, bindings)
}

func TestGroupByWithMultipleRawAndColumns(t *testing.T) {
	qb := xqb.Table("events")
	sql, _, _ := qb.
		Select("type", xqb.Raw("DATE(created_at) as day")).
		GroupBy("type", xqb.Raw("DATE(created_at)")).
		ToSQL()

	expected := "SELECT type, DATE(created_at) as day FROM events GROUP BY type, DATE(created_at)"
	assert.Equal(t, expected, sql)
}

func TestGroupByNoColumns(t *testing.T) {
	qb := xqb.Table("metrics")
	sql, _, _ := qb.
		Select("id").
		GroupBy().
		ToSQL()

	expected := "SELECT id FROM metrics"
	assert.Equal(t, expected, sql)
}

func TestGroupByDuplicateColumns(t *testing.T) {
	qb := xqb.Table("items")
	sql, _, _ := qb.
		Select("category_id").
		GroupBy("category_id", "category_id").
		ToSQL()

	expected := "SELECT category_id FROM items GROUP BY category_id, category_id"
	assert.Equal(t, expected, sql)
}

func TestGroupByWithOrderBy(t *testing.T) {
	qb := xqb.Table("sessions")
	sql, _, _ := qb.
		Select("user_id", xqb.Raw("COUNT(*) as count")).
		GroupBy("user_id").
		OrderBy("count", "desc").
		ToSQL()

	expected := "SELECT user_id, COUNT(*) as count FROM sessions GROUP BY user_id ORDER BY count desc"
	assert.Equal(t, expected, sql)
}
