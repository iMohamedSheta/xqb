package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/stretchr/testify/assert"
)

func Test_Pluck_ValueAndKey(t *testing.T) {
	qb := xqb.Table("users").Where("name", "LIKE", "%mohamed%")

	sql, bindings, err := qb.PluckSQL("name", "id")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "SELECT name, id FROM users WHERE name LIKE ?", sql)
	assert.Equal(t, []any{"%mohamed%"}, bindings)
	assert.Equal(t, []any{"name", "id"}, qb.GetData().Columns)
}

func Test_Pluck_Value(t *testing.T) {
	qb := xqb.Table("users").Where("name", "LIKE", "%mohamed%")
	sql, bindings, err := qb.PluckSQL("name", "")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "SELECT name FROM users WHERE name LIKE ?", sql)
	assert.Equal(t, []any{"%mohamed%"}, bindings)
	assert.Equal(t, []any{"name"}, qb.GetData().Columns)
}

func Test_Pluck_Key(t *testing.T) {
	qb := xqb.Table("users").Where("name", "LIKE", "%mohamed%")
	sql, bindings, err := qb.PluckSQL("", "id")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "SELECT id FROM users WHERE name LIKE ?", sql)
	assert.Equal(t, []any{"%mohamed%"}, bindings)
	assert.Equal(t, []any{"id"}, qb.GetData().Columns)
}

func Test_Pluck_NoData(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, err := qb.PluckSQL("name", "id")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "SELECT name, id FROM users", sql)
	assert.Equal(t, []any(nil), bindings)
	assert.Equal(t, []any{"name", "id"}, qb.GetData().Columns)
}

func Test_Pluck_NoData_NoKey_ReturnError(t *testing.T) {
	qb := xqb.Table("users")
	qb.Where("id", "=", 1)
	sql, bindings, err := qb.PluckSQL("", "")
	assert.Error(t, err)
	assert.Equal(t, 0, len(bindings))
	assert.Equal(t, "", sql)
}

func Test_Pluck_ComplexQuery(t *testing.T) {
	qb := xqb.Table("users").Where("name", "LIKE", "%mohamed%")
	qb.Select("id", "name").Where("id", "=", 1)
	sql, bindings, err := qb.PluckSQL("", "")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "SELECT id, name FROM users WHERE name LIKE ? AND id = ?", sql)
	assert.Equal(t, []any{"%mohamed%", 1}, bindings)
	assert.Equal(t, []any{"id", "name"}, qb.GetData().Columns)
}
