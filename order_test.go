package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/stretchr/testify/assert"
)

func TestOrderByWithRawExpressions(t *testing.T) {
	qb := xqb.Table("users")
	sql, _, _ := qb.OrderBy(xqb.Raw("FIELD(status, 'active', 'pending', 'inactive')"), "ASC").ToSQL()

	expected := "SELECT * FROM users ORDER BY FIELD(status, 'active', 'pending', 'inactive') ASC"
	assert.Equal(t, expected, sql)
}

func TestOrderBySimpleColumn(t *testing.T) {
	qb := xqb.Table("users").Select("*").OrderBy("name", "ASC")
	sql, _, _ := qb.ToSQL()

	expected := "SELECT * FROM users ORDER BY name ASC"
	assert.Equal(t, expected, sql)
}

func TestOrderByDescShortcut(t *testing.T) {
	qb := xqb.Table("users").Select("*").OrderByDesc("created_at")
	sql, _, _ := qb.ToSQL()

	expected := "SELECT * FROM users ORDER BY created_at DESC"
	assert.Equal(t, expected, sql)
}

func TestOrderByAscShortcut(t *testing.T) {
	qb := xqb.Table("users").Select("*").OrderByAsc("email")
	sql, _, _ := qb.ToSQL()

	expected := "SELECT * FROM users ORDER BY email ASC"
	assert.Equal(t, expected, sql)
}

func TestOrderByWithRawExpression(t *testing.T) {
	qb := xqb.Table("products").Select("*").
		OrderBy(xqb.Raw("LENGTH(name)"), "DESC")
	sql, bindings, _ := qb.ToSQL()

	expected := "SELECT * FROM products ORDER BY LENGTH(name) DESC"
	assert.Equal(t, expected, sql)
	assert.Empty(t, bindings)
}

func TestOrderByRawFunction(t *testing.T) {
	qb := xqb.Table("logs").Select("*").
		OrderByRaw("FIELD(status, ?, ?, ?)", "active", "pending", "disabled")
	sql, bindings, _ := qb.ToSQL()

	expected := "SELECT * FROM logs ORDER BY FIELD(status, ?, ?, ?)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{"active", "pending", "disabled"}, bindings)
}

func TestOrderByWithFallbackToString(t *testing.T) {
	qb := xqb.Table("items").Select("*").
		OrderBy(123, "ASC")
	sql, _, _ := qb.ToSQL()

	expected := "SELECT * FROM items ORDER BY 123 ASC"
	assert.Equal(t, expected, sql)
}

func TestLatestAndOldest(t *testing.T) {
	qb := xqb.Table("comments").Select("*").
		Latest("created_at").
		Oldest("updated_at")
	sql, _, _ := qb.ToSQL()

	expected := "SELECT * FROM comments ORDER BY created_at DESC, updated_at ASC"
	assert.Equal(t, expected, sql)
}
