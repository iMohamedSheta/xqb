package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/stretchr/testify/assert"
)

func TestOrderByWithRawExpressions(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.OrderBy(xqb.Raw("FIELD(status, 'active', 'pending', 'inactive')"), "ASC").ToSql()
		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "SELECT * FROM `users` ORDER BY FIELD(status, 'active', 'pending', 'inactive') ASC",
			types.DialectPostgres: `SELECT * FROM "users" ORDER BY FIELD(status, 'active', 'pending', 'inactive') ASC`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func TestOrderBySimpleColumn(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Select("*").OrderBy("name", "ASC")
		sql, bindings, err := qb.ToSql()

		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "SELECT * FROM `users` ORDER BY `name` ASC",
			types.DialectPostgres: `SELECT * FROM "users" ORDER BY "name" ASC`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func TestOrderByDescShortcut(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Select("*").OrderByDesc("created_at")
		sql, bindings, err := qb.ToSql()

		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "SELECT * FROM `users` ORDER BY `created_at` DESC",
			types.DialectPostgres: `SELECT * FROM "users" ORDER BY "created_at" DESC`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func TestOrderByAscShortcut(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Select("*").OrderByAsc("email")
		sql, bindings, err := qb.ToSql()

		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "SELECT * FROM `users` ORDER BY `email` ASC",
			types.DialectPostgres: `SELECT * FROM "users" ORDER BY "email" ASC`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func TestOrderByWithRawExpression(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("products").Select("*").SetDialect(dialect).
			OrderBy(xqb.Raw("LENGTH(name)"), "DESC")
		sql, bindings, err := qb.ToSql()

		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "SELECT * FROM `products` ORDER BY LENGTH(name) DESC",
			types.DialectPostgres: `SELECT * FROM "products" ORDER BY LENGTH(name) DESC`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func TestOrderByRawFunction(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("logs").SetDialect(dialect).Select("*").
			OrderByRaw("FIELD(status, ?, ?, ?)", "active", "pending", "disabled")
		sql, bindings, err := qb.ToSql()

		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "SELECT * FROM `logs` ORDER BY FIELD(status, ?, ?, ?)",
			types.DialectPostgres: `SELECT * FROM "logs" ORDER BY FIELD(status, $1, $2, $3)`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"active", "pending", "disabled"}, bindings)
		assert.NoError(t, err)
	})
}

func TestOrderByWithFallbackToString(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("items").SetDialect(dialect).
			Select("*").
			OrderBy(123, "ASC")

		sql, bindings, err := qb.ToSql()

		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "SELECT * FROM `items` ORDER BY 123 ASC",
			types.DialectPostgres: `SELECT * FROM "items" ORDER BY 123 ASC`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func TestLatestAndOldest(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("comments").SetDialect(dialect).
			Select("*").
			Latest("created_at").
			Oldest("updated_at")

		sql, bindings, err := qb.ToSql()

		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "SELECT * FROM `comments` ORDER BY `created_at` DESC, `updated_at` ASC",
			types.DialectPostgres: `SELECT * FROM "comments" ORDER BY "created_at" DESC, "updated_at" ASC`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}
