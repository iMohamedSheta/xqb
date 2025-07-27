package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/stretchr/testify/assert"
)

func Test_Union(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		q := xqb.Table("users").SetDialect(dialect).
			Select("id", "name").
			UnionRaw("SELECT id, name FROM admins WHERE active = ?", true).
			Union(
				xqb.Table("admins").SetDialect(dialect).
					Select("id", "username").
					Where("username", "=", "mohamed").
					Limit(1),
			)

		sql, bindings, err := q.ToSql()
		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "(SELECT `id`, `name` FROM `users`) UNION (SELECT id, name FROM admins WHERE active = ?) UNION (SELECT `id`, `username` FROM `admins` WHERE `username` = ? LIMIT 1)",
			types.DialectPostgres: `(SELECT "id", "name" FROM "users") UNION (SELECT id, name FROM admins WHERE active = $1) UNION (SELECT "id", "username" FROM "admins" WHERE "username" = $2 LIMIT 1)`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		expectedBindings := []any{true, "mohamed"}
		assert.Equal(t, expectedBindings, bindings)
		assert.NoError(t, err)
	})
}

func Test_UnionAll(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		q := xqb.Table("users").SetDialect(dialect).
			Select("id").
			UnionAllRaw("SELECT id FROM guests WHERE banned = ?", false)

		sql, bindings, err := q.ToSql()
		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "(SELECT `id` FROM `users`) UNION ALL (SELECT id FROM guests WHERE banned = ?)",
			types.DialectPostgres: `(SELECT "id" FROM "users") UNION ALL (SELECT id FROM guests WHERE banned = $1)`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		expectedBindings := []any{false}
		assert.Equal(t, expectedBindings, bindings)
		assert.NoError(t, err)
	})
}

func Test_ExceptUnion(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		q := xqb.Table("users").SetDialect(dialect).
			Select("id").
			ExceptUnionRaw("SELECT id FROM banned_users", false)

		sql, bindings, err := q.ToSql()
		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "", // not Supported by MySql
			types.DialectPostgres: `(SELECT "id" FROM "users") EXCEPT (SELECT id FROM banned_users)`,
		}
		expectedErr := map[types.Dialect]error{
			types.DialectMySql:    errors.ErrUnsupportedFeature,
			types.DialectPostgres: nil,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		if expectedErr[dialect] != nil {
			assert.ErrorIs(t, err, expectedErr[dialect])
		} else {
			assert.NoError(t, err)
		}
		assert.Empty(t, bindings)
	})
}

func Test_ExceptUnion_All(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		q := xqb.Table("users").SetDialect(dialect).
			Select("id").
			ExceptUnionRaw("SELECT id FROM banned_users", true)

		sql, bindings, err := q.ToSql()
		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "", // not Supported by MySql
			types.DialectPostgres: `(SELECT "id" FROM "users") EXCEPT ALL (SELECT id FROM banned_users)`,
		}
		expectedErr := map[types.Dialect]error{
			types.DialectMySql:    errors.ErrUnsupportedFeature,
			types.DialectPostgres: nil,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		if expectedErr[dialect] != nil {
			assert.ErrorIs(t, err, expectedErr[dialect])
		} else {
			assert.NoError(t, err)
		}
		assert.Empty(t, bindings)
	})
}

func Test_IntersectUnion(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		q := xqb.Table("users").SetDialect(dialect).
			Select("id").
			IntersectUnionRaw("SELECT id FROM employees WHERE active = ?", true, true)

		sql, bindings, err := q.ToSql()
		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "", // not Supported by MySql
			types.DialectPostgres: `(SELECT "id" FROM "users") INTERSECT ALL (SELECT id FROM employees WHERE active = $1)`,
		}
		expectedErr := map[types.Dialect]error{
			types.DialectMySql:    errors.ErrUnsupportedFeature,
			types.DialectPostgres: nil,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		if expectedErr[dialect] != nil {
			assert.ErrorIs(t, err, expectedErr[dialect])
		} else {
			expectedBindings := []any{true}
			assert.Equal(t, expectedBindings, bindings)
			assert.NoError(t, err)
		}
	})
}

func Test_Union_WithMultipleQueries(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		q := xqb.Table("users").SetDialect(dialect).Select("id").
			Union(
				xqb.Table("admins").SetDialect(dialect).Select("id"),
				xqb.Table("guests").SetDialect(dialect).Select("id"),
			)

		sql, bindings, err := q.ToSql()
		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "(SELECT `id` FROM `users`) UNION (SELECT `id` FROM `admins`) UNION (SELECT `id` FROM `guests`)",
			types.DialectPostgres: `(SELECT "id" FROM "users") UNION (SELECT "id" FROM "admins") UNION (SELECT "id" FROM "guests")`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_UnionAll_WithMultipleQueries(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		q := xqb.Table("users").SetDialect(dialect).Select("id").
			UnionAll(
				xqb.Table("admins").SetDialect(dialect).Select("id"),
				xqb.Table("guests").SetDialect(dialect).Select("id"),
			)

		sql, bindings, err := q.ToSql()
		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "(SELECT `id` FROM `users`) UNION ALL (SELECT `id` FROM `admins`) UNION ALL (SELECT `id` FROM `guests`)",
			types.DialectPostgres: `(SELECT "id" FROM "users") UNION ALL (SELECT "id" FROM "admins") UNION ALL (SELECT "id" FROM "guests")`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_Union_MixedRawAndBuilder(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		q := xqb.Table("users").SetDialect(dialect).Select("id").
			UnionRaw("SELECT id FROM guests WHERE active = ?", true).
			Union(xqb.Table("admins").SetDialect(dialect).Select("id").Where("id", ">", 5))

		sql, bindings, err := q.ToSql()
		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "(SELECT `id` FROM `users`) UNION (SELECT id FROM guests WHERE active = ?) UNION (SELECT `id` FROM `admins` WHERE `id` > ?)",
			types.DialectPostgres: `(SELECT "id" FROM "users") UNION (SELECT id FROM guests WHERE active = $1) UNION (SELECT "id" FROM "admins" WHERE "id" > $2)`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		expectedBindings := []any{true, 5}
		assert.Equal(t, expectedBindings, bindings)
		assert.NoError(t, err)
	})
}

func Test_UnionAllRaw_WithBindings(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		q := xqb.Table("users").SetDialect(dialect).Select("id").
			UnionAllRaw("SELECT id FROM banned_users WHERE reason = ?", "spam")

		sql, bindings, err := q.ToSql()
		assert.NoError(t, err)
		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "(SELECT `id` FROM `users`) UNION ALL (SELECT id FROM banned_users WHERE reason = ?)",
			types.DialectPostgres: `(SELECT "id" FROM "users") UNION ALL (SELECT id FROM banned_users WHERE reason = $1)`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		expectedBindings := []any{"spam"}
		assert.Equal(t, expectedBindings, bindings)
		assert.NoError(t, err)
	})
}

func Test_ExceptUnion_Unsupported(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		q := xqb.Table("users").SetDialect(dialect).Select("id").
			ExceptUnion(xqb.Table("banned_users").SetDialect(dialect).Select("id"))

		sql, bindings, err := q.ToSql()
		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "", // not Supported by MySql
			types.DialectPostgres: `(SELECT "id" FROM "users") EXCEPT (SELECT "id" FROM "banned_users")`,
		}
		expectedErr := map[types.Dialect]error{
			types.DialectMySql:    errors.ErrUnsupportedFeature,
			types.DialectPostgres: nil,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		if expectedErr[dialect] != nil {
			assert.ErrorIs(t, err, expectedErr[dialect])
		} else {
			assert.NoError(t, err)
		}
		assert.Empty(t, bindings)
	})
}

func Test_ExceptUnionAll_Unsupported(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		q := xqb.Table("users").SetDialect(dialect).Select("id").
			ExceptUnionAll(xqb.Table("banned_users").SetDialect(dialect).Select("id"))

		sql, bindings, err := q.ToSql()

		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "", // not Supported by MySql
			types.DialectPostgres: `(SELECT "id" FROM "users") EXCEPT ALL (SELECT "id" FROM "banned_users")`,
		}

		expectedErr := map[types.Dialect]error{
			types.DialectMySql:    errors.ErrUnsupportedFeature,
			types.DialectPostgres: nil,
		}
		assert.Equal(t, expectedSql[dialect], sql)

		if expectedErr[dialect] != nil {
			assert.ErrorIs(t, err, expectedErr[dialect])
		} else {
			assert.NoError(t, err)
		}

		assert.Empty(t, bindings)
	})
}

func Test_ExceptUnionRaw_Unsupported(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		q := xqb.Table("users").SetDialect(dialect).Select("id").
			ExceptUnionRaw("SELECT id FROM banned_users", true)

		sql, bindings, err := q.ToSql()

		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "", // not Supported by MySql
			types.DialectPostgres: `(SELECT "id" FROM "users") EXCEPT ALL (SELECT id FROM banned_users)`,
		}

		expectedErr := map[types.Dialect]error{
			types.DialectMySql:    errors.ErrUnsupportedFeature,
			types.DialectPostgres: nil,
		}

		assert.Equal(t, expectedSql[dialect], sql)

		if expectedErr[dialect] != nil {
			assert.ErrorIs(t, err, expectedErr[dialect])
		} else {
			assert.NoError(t, err)
		}

		assert.Empty(t, bindings)
	})
}

func Test_IntersectUnion_Unsupported(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		q := xqb.Table("users").SetDialect(dialect).Select("id").
			IntersectUnion(xqb.Table("employees").SetDialect(dialect).Select("id").Where("active", "=", true))

		sql, bindings, err := q.ToSql()
		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "", // not Supported by MySql
			types.DialectPostgres: `(SELECT "id" FROM "users") INTERSECT (SELECT "id" FROM "employees" WHERE "active" = $1)`,
		}
		expectedErr := map[types.Dialect]error{
			types.DialectMySql:    errors.ErrUnsupportedFeature,
			types.DialectPostgres: nil,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		if expectedErr[dialect] != nil {
			assert.ErrorIs(t, err, expectedErr[dialect])
		} else {
			expectedBindings := []any{true}
			assert.Equal(t, expectedBindings, bindings)
			assert.NoError(t, err)
		}
	})
}

func Test_IntersectUnionAll_Unsupported(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		q := xqb.Table("users").SetDialect(dialect).Select("id").
			IntersectUnionAll(xqb.Table("employees").SetDialect(dialect).Select("id").Where("active", "=", true))

		sql, bindings, err := q.ToSql()
		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "", // not Supported by MySql
			types.DialectPostgres: `(SELECT "id" FROM "users") INTERSECT ALL (SELECT "id" FROM "employees" WHERE "active" = $1)`,
		}
		expectedErr := map[types.Dialect]error{
			types.DialectMySql:    errors.ErrUnsupportedFeature,
			types.DialectPostgres: nil,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		if expectedErr[dialect] != nil {
			assert.ErrorIs(t, err, expectedErr[dialect])
		} else {
			expectedBindings := []any{true}
			assert.Equal(t, expectedBindings, bindings)
			assert.NoError(t, err)
		}
	})
}

func Test_IntersectUnionRaw_Unsupported(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		q := xqb.Table("users").SetDialect(dialect).Select("id").
			IntersectUnionRaw("SELECT id FROM employees", false)

		sql, bindings, err := q.ToSql()
		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "", // not Supported by MySql
			types.DialectPostgres: `(SELECT "id" FROM "users") INTERSECT (SELECT id FROM employees)`,
		}
		expectedErr := map[types.Dialect]error{
			types.DialectMySql:    errors.ErrUnsupportedFeature,
			types.DialectPostgres: nil,
		}
		assert.Equal(t, expectedSql[dialect], sql)

		if expectedErr[dialect] != nil {
			assert.ErrorIs(t, err, expectedErr[dialect])
		} else {
			assert.NoError(t, err)
		}

		assert.Empty(t, bindings)
	})
}

func Test_Union_WithEmptyUnionList(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		q := xqb.Table("users").SetDialect(dialect).Select("id")
		sql, bindings, err := q.ToSql()
		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "SELECT `id` FROM `users`",
			types.DialectPostgres: `SELECT "id" FROM "users"`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}
