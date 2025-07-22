package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/stretchr/testify/assert"
)

func Test_Having_WithRawExpressions(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("orders").SetDialect(dialect)
		sql, bindings, err := qb.Select("user_id", xqb.Raw("SUM(amount) as total")).
			GroupBy("user_id").
			Having(xqb.Raw("SUM(amount)"), ">", 1000).
			ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT user_id, SUM(amount) as total FROM orders GROUP BY user_id HAVING SUM(amount) > ?",
			types.DriverPostgres: "SELECT user_id, SUM(amount) as total FROM orders GROUP BY user_id HAVING SUM(amount) > $1",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{1000}, bindings)
		assert.NoError(t, err)
	})
}

func Test_Having_Simple(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("orders").SetDialect(dialect)
		sql, bindings, err := qb.Select("user_id", xqb.Raw("SUM(amount) as total")).
			GroupBy("user_id").
			Having("SUM(amount)", ">", 500).
			ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT user_id, SUM(amount) as total FROM orders GROUP BY user_id HAVING SUM(amount) > ?",
			types.DriverPostgres: "SELECT user_id, SUM(amount) as total FROM orders GROUP BY user_id HAVING SUM(amount) > $1",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{500}, bindings)
		assert.NoError(t, err)

	})
}

func Test_Having_Raw(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("orders").SetDialect(dialect)
		sql, bindings, err := qb.Select("user_id", xqb.Raw("SUM(amount) as total")).
			GroupBy("user_id").
			HavingRaw("SUM(amount) > ?", 1000).
			ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT user_id, SUM(amount) as total FROM orders GROUP BY user_id HAVING SUM(amount) > ?",
			types.DriverPostgres: "SELECT user_id, SUM(amount) as total FROM orders GROUP BY user_id HAVING SUM(amount) > $1",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{1000}, bindings)
		assert.NoError(t, err)
	})
}

func Test_OrHaving_WithExpressions(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("orders").SetDialect(dialect)
		sql, bindings, err := qb.Select("user_id", xqb.Raw("SUM(amount) as total")).
			GroupBy("user_id").
			Having("SUM(amount)", ">", 1000).
			OrHaving("SUM(discount)", ">", 200).
			ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT user_id, SUM(amount) as total FROM orders GROUP BY user_id HAVING SUM(amount) > ? OR SUM(discount) > ?",
			types.DriverPostgres: "SELECT user_id, SUM(amount) as total FROM orders GROUP BY user_id HAVING SUM(amount) > $1 OR SUM(discount) > $2",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{1000, 200}, bindings)
		assert.NoError(t, err)
	})
}

func Test_OrHaving_WithRaw(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("orders").SetDialect(dialect)
		sql, bindings, err := qb.Select("user_id", xqb.Raw("SUM(amount) as total")).
			GroupBy("user_id").
			HavingRaw("SUM(amount) > ?", 1000).
			OrHavingRaw("SUM(discount) > ?", 200).
			ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT user_id, SUM(amount) as total FROM orders GROUP BY user_id HAVING SUM(amount) > ? OR SUM(discount) > ?",
			types.DriverPostgres: "SELECT user_id, SUM(amount) as total FROM orders GROUP BY user_id HAVING SUM(amount) > $1 OR SUM(discount) > $2",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{1000, 200}, bindings)
		assert.NoError(t, err)
	})
}

func Test_Having_WithMultiple(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("orders").SetDialect(dialect)
		sql, bindings, err := qb.Select("user_id", xqb.Raw("SUM(amount) as total")).
			GroupBy("user_id").
			Having("SUM(amount)", ">", 1000).
			Having("SUM(discount)", "<", 100).
			ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT user_id, SUM(amount) as total FROM orders GROUP BY user_id HAVING SUM(amount) > ? AND SUM(discount) < ?",
			types.DriverPostgres: "SELECT user_id, SUM(amount) as total FROM orders GROUP BY user_id HAVING SUM(amount) > $1 AND SUM(discount) < $2",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{1000, 100}, bindings)
		assert.NoError(t, err)
	})
}

func Test_Having_WithExpressionToExpression(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		expr1 := xqb.Raw("SUM(amount)")
		expr2 := xqb.Raw("SUM(discount)")
		qb := xqb.Table("orders").SetDialect(dialect)

		sql, bindings, err := qb.Select("user_id", xqb.Raw("SUM(amount) as total")).
			GroupBy("user_id").
			Having(expr1, ">", expr2).
			ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT user_id, SUM(amount) as total FROM orders GROUP BY user_id HAVING SUM(amount) > SUM(discount)",
			types.DriverPostgres: "SELECT user_id, SUM(amount) as total FROM orders GROUP BY user_id HAVING SUM(amount) > SUM(discount)",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_Having_WithIsNull(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("orders").SetDialect(dialect)
		sql, bindings, err := qb.Select("user_id", xqb.Raw("COUNT(*) as cnt")).
			GroupBy("user_id").
			Having("COUNT(*)", "IS NULL", nil).
			ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT user_id, COUNT(*) as cnt FROM orders GROUP BY user_id HAVING COUNT(*) IS NULL",
			types.DriverPostgres: "SELECT user_id, COUNT(*) as cnt FROM orders GROUP BY user_id HAVING COUNT(*) IS NULL",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_Having_WithExpressionValue(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("orders").SetDialect(dialect)
		sql, bindings, err := qb.Select("user_id").
			GroupBy("user_id").
			Having(xqb.Raw("SUM(amount)"), ">", xqb.Raw("AVG(amount)")).
			ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT user_id FROM orders GROUP BY user_id HAVING SUM(amount) > AVG(amount)",
			types.DriverPostgres: "SELECT user_id FROM orders GROUP BY user_id HAVING SUM(amount) > AVG(amount)",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_Having_WithExpressionAndBoundValue(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("orders").SetDialect(dialect)
		sql, bindings, err := qb.Select("user_id").
			GroupBy("user_id").
			Having(xqb.Raw("SUM(amount)"), ">", 100).
			ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT user_id FROM orders GROUP BY user_id HAVING SUM(amount) > ?",
			types.DriverPostgres: "SELECT user_id FROM orders GROUP BY user_id HAVING SUM(amount) > $1",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{100}, bindings)
		assert.NoError(t, err)
	})
}

func Test_Having_WithExpressionValueAndBindings(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("orders").SetDialect(dialect)
		sql, bindings, err := qb.Select("user_id").
			GroupBy("user_id").
			Having(xqb.Raw("COALESCE(SUM(amount), ?)", 15), "=", xqb.Raw("?", 25)).
			ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT user_id FROM orders GROUP BY user_id HAVING COALESCE(SUM(amount), ?) = ?",
			types.DriverPostgres: "SELECT user_id FROM orders GROUP BY user_id HAVING COALESCE(SUM(amount), $1) = $2",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{15, 25}, bindings)
		assert.NoError(t, err)
	})
}
