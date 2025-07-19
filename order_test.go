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
