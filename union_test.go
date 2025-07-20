package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/stretchr/testify/assert"
)

func Test_Union(t *testing.T) {
	q := xqb.Table("users").
		Select("id", "name").
		UnionRaw("SELECT id, name FROM admins WHERE active = ?", true)

	sql, bindings, err := q.ToSQL()
	assert.NoError(t, err)

	expected := "SELECT id, name FROM users UNION (SELECT id, name FROM admins WHERE active = ?)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{true}, bindings)
}

func Test_UnionAll(t *testing.T) {
	q := xqb.Table("users").
		Select("id").
		UnionAllRaw("SELECT id FROM guests WHERE banned = ?", false)

	sql, bindings, err := q.ToSQL()
	assert.NoError(t, err)

	expected := "SELECT id FROM users UNION ALL (SELECT id FROM guests WHERE banned = ?)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{false}, bindings)
}

func Test_ExceptUnion(t *testing.T) {
	q := xqb.Table("users").
		Select("id").
		ExceptUnionRaw("SELECT id FROM banned_users", false)

	_, _, err := q.ToSQL()
	assert.ErrorContains(t, err, "union type Except is not supported in MySQL")
}

func Test_ExceptUnion_All(t *testing.T) {
	q := xqb.Table("users").
		Select("id").
		ExceptUnionRaw("SELECT id FROM banned_users", true)

	_, _, err := q.ToSQL()
	assert.ErrorContains(t, err, "union type Except is not supported in MySQL")
}

func Test_IntersectUnion(t *testing.T) {
	q := xqb.Table("users").
		Select("id").
		IntersectUnionRaw("SELECT id FROM employees WHERE active = ?", true, true)

	_, _, err := q.ToSQL()
	assert.ErrorContains(t, err, "union type Intersect is not supported in MySQL")
}
