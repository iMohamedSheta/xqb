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
