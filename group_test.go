package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/stretchr/testify/assert"
)

func TestGroupByWithRawExpressions(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("orders").SetDialect(dialect)
		sql, bindings, err := qb.Select(xqb.Raw("YEAR(created_at) as year"), xqb.Raw("SUM(amount) as total")).
			GroupBy(xqb.Raw("YEAR(created_at)")).
			ToSql()

		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "SELECT YEAR(created_at) as year, SUM(amount) as total FROM `orders` GROUP BY YEAR(created_at)",
			types.DialectPostgres: `SELECT YEAR(created_at) as year, SUM(amount) as total FROM "orders" GROUP BY YEAR(created_at)`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func TestGroupByMultipleColumns(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("orders").SetDialect(dialect)
		sql, bindings, err := qb.Select("user_id", "product_id").
			GroupBy("user_id", "product_id").
			ToSql()

		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "SELECT `user_id`, `product_id` FROM `orders` GROUP BY `user_id`, `product_id`",
			types.DialectPostgres: `SELECT "user_id", "product_id" FROM "orders" GROUP BY "user_id", "product_id"`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func TestGroupByRawShortcut(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("orders").SetDialect(dialect)
		sql, bindings, err := qb.Select("id").
			GroupByRaw("DATE(created_at)").
			ToSql()

		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "SELECT `id` FROM `orders` GROUP BY DATE(created_at)",
			types.DialectPostgres: `SELECT "id" FROM "orders" GROUP BY DATE(created_at)`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func TestGroupByWithHaving(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("sales").SetDialect(dialect)
		sql, bindings, err := qb.
			Select("region", xqb.Raw("SUM(amount) as total")).
			GroupBy("region").
			Having("SUM(amount)", ">", 1000).
			ToSql()

		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "SELECT `region`, SUM(amount) as total FROM `sales` GROUP BY `region` HAVING SUM(amount) > ?",
			types.DialectPostgres: `SELECT "region", SUM(amount) as total FROM "sales" GROUP BY "region" HAVING SUM(amount) > $1`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{1000}, bindings)
		assert.NoError(t, err)
	})
}

func TestGroupByWithMultipleRawAndColumns(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("events").SetDialect(dialect)
		sql, bindings, err := qb.
			Select("type", xqb.Raw("DATE(created_at) as day")).
			GroupBy("type", xqb.Raw("DATE(created_at)")).
			ToSql()

		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "SELECT `type`, DATE(created_at) as day FROM `events` GROUP BY `type`, DATE(created_at)",
			types.DialectPostgres: `SELECT "type", DATE(created_at) as day FROM "events" GROUP BY "type", DATE(created_at)`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func TestGroupByNoColumns_ReturnAnError(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("metrics").SetDialect(dialect)
		sql, bindings, err := qb.
			Select("id").
			GroupBy().
			ToSql()

		assert.ErrorIs(t, err, errors.ErrInvalidQuery)
		assert.Empty(t, sql)
		assert.Empty(t, bindings)
	})
}

func TestGroupByWithOrderBy(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("sessions").SetDialect(dialect)
		sql, bindings, err := qb.
			Select("user_id", xqb.Raw("COUNT(*) as count")).
			GroupBy("user_id").
			OrderBy("count", "DESC").
			ToSql()

		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "SELECT `user_id`, COUNT(*) as count FROM `sessions` GROUP BY `user_id` ORDER BY `count` DESC",
			types.DialectPostgres: `SELECT "user_id", COUNT(*) as count FROM "sessions" GROUP BY "user_id" ORDER BY "count" DESC`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}
