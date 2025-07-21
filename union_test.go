package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/stretchr/testify/assert"
)

func Test_Union(t *testing.T) {
	q := xqb.Table("users").
		Select("id", "name").
		UnionRaw("SELECT id, name FROM admins WHERE active = ?", true).
		Union(
			xqb.Table("admins").
				Select("id", "username").
				Where("username", "=", "mohamed").
				Limit(1),
		)

	sql, bindings, err := q.ToSQL()
	assert.NoError(t, err)

	expected := "(SELECT id, name FROM users) UNION (SELECT id, name FROM admins WHERE active = ?) UNION (SELECT id, username FROM admins WHERE username = ? LIMIT 1)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{true, "mohamed"}, bindings)
}

func Test_UnionAll(t *testing.T) {
	q := xqb.Table("users").
		Select("id").
		UnionAllRaw("SELECT id FROM guests WHERE banned = ?", false)

	sql, bindings, err := q.ToSQL()
	assert.NoError(t, err)

	expected := "(SELECT id FROM users) UNION ALL (SELECT id FROM guests WHERE banned = ?)"
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

func Test_Union_WithMultipleQueries(t *testing.T) {
	q := xqb.Table("users").Select("id").
		Union(
			xqb.Table("admins").Select("id"),
			xqb.Table("guests").Select("id"),
		)

	sql, _, err := q.ToSQL()
	assert.NoError(t, err)
	expected := "(SELECT id FROM users) UNION (SELECT id FROM admins) UNION (SELECT id FROM guests)"
	assert.Equal(t, expected, sql)
}

func Test_UnionAll_WithMultipleQueries(t *testing.T) {
	q := xqb.Table("users").Select("id").
		UnionAll(
			xqb.Table("admins").Select("id"),
			xqb.Table("guests").Select("id"),
		)

	sql, _, err := q.ToSQL()
	assert.NoError(t, err)
	expected := "(SELECT id FROM users) UNION ALL (SELECT id FROM admins) UNION ALL (SELECT id FROM guests)"
	assert.Equal(t, expected, sql)
}

func Test_Union_MixedRawAndBuilder(t *testing.T) {
	q := xqb.Table("users").Select("id").
		UnionRaw("SELECT id FROM guests WHERE active = ?", true).
		Union(xqb.Table("admins").Select("id").Where("id", ">", 5))

	sql, bindings, err := q.ToSQL()
	assert.NoError(t, err)
	expected := "(SELECT id FROM users) UNION (SELECT id FROM guests WHERE active = ?) UNION (SELECT id FROM admins WHERE id > ?)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{true, 5}, bindings)
}

func Test_UnionAllRaw_WithBindings(t *testing.T) {
	q := xqb.Table("users").Select("id").
		UnionAllRaw("SELECT id FROM banned_users WHERE reason = ?", "spam")

	sql, bindings, err := q.ToSQL()
	assert.NoError(t, err)
	expected := "(SELECT id FROM users) UNION ALL (SELECT id FROM banned_users WHERE reason = ?)"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{"spam"}, bindings)
}

func Test_ExceptUnion_Unsupported(t *testing.T) {
	q := xqb.Table("users").Select("id").
		ExceptUnion(xqb.Table("banned_users").Select("id"))

	_, _, err := q.ToSQL()
	assert.ErrorContains(t, err, "union type Except is not supported in MySQL")
}

func Test_ExceptUnionAll_Unsupported(t *testing.T) {
	q := xqb.Table("users").Select("id").
		ExceptUnionAll(xqb.Table("banned_users").Select("id"))

	_, _, err := q.ToSQL()
	assert.ErrorContains(t, err, "union type Except is not supported in MySQL")
}

func Test_ExceptUnionRaw_Unsupported(t *testing.T) {
	q := xqb.Table("users").Select("id").
		ExceptUnionRaw("SELECT id FROM banned_users", true)

	_, _, err := q.ToSQL()
	assert.ErrorContains(t, err, "union type Except is not supported in MySQL")
}

func Test_IntersectUnion_Unsupported(t *testing.T) {
	q := xqb.Table("users").Select("id").
		IntersectUnion(xqb.Table("employees").Select("id"))

	_, _, err := q.ToSQL()
	assert.ErrorContains(t, err, "union type Intersect is not supported in MySQL")
}

func Test_IntersectUnionAll_Unsupported(t *testing.T) {
	q := xqb.Table("users").Select("id").
		IntersectUnionAll(xqb.Table("employees").Select("id"))

	_, _, err := q.ToSQL()
	assert.ErrorContains(t, err, "union type Intersect is not supported in MySQL")
}

func Test_IntersectUnionRaw_Unsupported(t *testing.T) {
	q := xqb.Table("users").Select("id").
		IntersectUnionRaw("SELECT id FROM employees", false)

	_, _, err := q.ToSQL()
	assert.ErrorContains(t, err, "union type Intersect is not supported in MySQL")
}

func Test_Union_WithEmptyUnionList(t *testing.T) {
	q := xqb.Table("users").Select("id")

	sql, _, err := q.ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT id FROM users", sql)
}
