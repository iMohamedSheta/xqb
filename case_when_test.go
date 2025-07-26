package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/stretchr/testify/assert"
)

func Test_CaseWhen(t *testing.T) {
	caseBuilder := xqb.Case()
	caseBuilder.When("age > ?", "adult", 18)
	caseBuilder.When("age < ?", "minor", 18)
	caseBuilder.When("age > ?", "dead", 100)
	caseBuilder.Else("dead")
	sql, bindings, err := caseBuilder.End().ToSql()
	assert.Equal(t, "CASE WHEN age > ? THEN ? WHEN age < ? THEN ? WHEN age > ? THEN ? ELSE ? END", sql)
	assert.Equal(t, []any{18, "adult", 18, "minor", 100, "dead", "dead"}, bindings)
	assert.NoError(t, err)
}

func Test_CaseWhen_UsageInQuery(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		caseExpr := xqb.Case().
			When("age >= ?", "adult", 18).
			When("age < ?", "minor", 18).
			Else("unknown").
			As("age_group").
			End()
		sql, bindings, err := qb.Select("id", caseExpr).
			Where(caseExpr, "=", "adult").
			Having(xqb.Count("id", ""), ">", 10).
			ToSql()

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT `id`, CASE WHEN age >= ? THEN ? WHEN age < ? THEN ? ELSE ? END AS age_group FROM `users` WHERE CASE WHEN age >= ? THEN ? WHEN age < ? THEN ? ELSE ? END AS age_group = ? HAVING COUNT(id) > ?",
			types.DriverPostgres: `SELECT "id", CASE WHEN age >= $1 THEN $2 WHEN age < $3 THEN $4 ELSE $5 END AS age_group FROM "users" WHERE CASE WHEN age >= $6 THEN $7 WHEN age < $8 THEN $9 ELSE $10 END AS age_group = $11 HAVING COUNT(id) > $12`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{18, "adult", 18, "minor", "unknown", 18, "adult", 18, "minor", "unknown", "adult", 10}, bindings)
		assert.NoError(t, err)
	})
}

func Test_EmptyAliasAndNoBindings(t *testing.T) {
	expr := xqb.Sum("amount", "")
	assert.Equal(t, "SUM(amount)", expr.Sql)
	expr2 := xqb.Case().When("1=1", "yes").End()
	assert.Equal(t, "CASE WHEN 1=1 THEN ? END", expr2.Sql)
}

func Test_CaseWhen_SingleWhen(t *testing.T) {
	caseBuilder := xqb.Case()
	caseBuilder.When("score > ?", "pass", 50)
	sql, bindings, err := caseBuilder.End().ToSql()
	expectedSql := "CASE WHEN score > ? THEN ? END"
	assert.Equal(t, expectedSql, sql)
	assert.Equal(t, []any{50, "pass"}, bindings)
	assert.NoError(t, err)
}

func Test_CaseWhen_MultipleWhen_NoElse(t *testing.T) {
	caseBuilder := xqb.Case()
	caseBuilder.When("score > ?", "A", 90)
	caseBuilder.When("score > ?", "B", 80)
	sql, bindings, err := caseBuilder.End().ToSql()
	expectedSql := "CASE WHEN score > ? THEN ? WHEN score > ? THEN ? END"
	assert.Equal(t, expectedSql, sql)
	assert.Equal(t, []any{90, "A", 80, "B"}, bindings)
	assert.NoError(t, err)
}

func Test_CaseWhen_ElseOnly(t *testing.T) {
	caseBuilder := xqb.Case()
	caseBuilder.Else("fail")
	sql, bindings, err := caseBuilder.End().ToSql()
	expectedSql := "CASE ELSE ? END"
	assert.Equal(t, expectedSql, sql)
	assert.Equal(t, []any{"fail"}, bindings)
	assert.NoError(t, err)
}

func Test_CaseWhen_WithAlias(t *testing.T) {
	caseBuilder := xqb.Case()
	caseBuilder.When("x = ?", "yes", 1).Else("no").As("result")
	sql, bindings, err := caseBuilder.End().ToSql()
	expectedSql := "CASE WHEN x = ? THEN ? ELSE ? END AS result"
	assert.Equal(t, expectedSql, sql)
	assert.Equal(t, []any{1, "yes", "no"}, bindings)
	assert.NoError(t, err)
}

func Test_CaseWhen_NoAlias(t *testing.T) {
	caseBuilder := xqb.Case()
	caseBuilder.When("x = ?", "yes", 1).Else("no")
	sql, bindings, err := caseBuilder.End().ToSql()
	expectedSql := "CASE WHEN x = ? THEN ? ELSE ? END"
	assert.Equal(t, expectedSql, sql)
	assert.Equal(t, []any{1, "yes", "no"}, bindings)
	assert.NoError(t, err)
}

func Test_CaseWhen_NoBindings(t *testing.T) {
	caseBuilder := xqb.Case()
	caseBuilder.When("1=1", "ok")
	sql, bindings, err := caseBuilder.End().ToSql()
	expectedSql := "CASE WHEN 1=1 THEN ? END"
	assert.Equal(t, expectedSql, sql)
	assert.Equal(t, []any{"ok"}, bindings)
	assert.NoError(t, err)
}

func Test_CaseWhen_ComplexConditions(t *testing.T) {
	caseBuilder := xqb.Case()
	caseBuilder.When("score > ? AND passed = ?", "excellent", 95, true)
	caseBuilder.When("score > ?", "good", 80)
	caseBuilder.Else("average")
	sql, bindings, err := caseBuilder.End().ToSql()
	expectedSql := "CASE WHEN score > ? AND passed = ? THEN ? WHEN score > ? THEN ? ELSE ? END"
	assert.Equal(t, expectedSql, sql)
	assert.Equal(t, []any{95, true, "excellent", 80, "good", "average"}, bindings)
	assert.NoError(t, err)
}

func Test_CaseWhen_SelectWithConditionalExpressions(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("orders").SetDialect(dialect)
		qb.Select(
			"id",
			xqb.Case().
				When("status = ?", "'done'", "'completed'").
				Else("'pending'").
				As("status_text").
				End(),
		)

		sql, bindings, err := qb.ToSql()

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT `id`, CASE WHEN status = ? THEN ? ELSE ? END AS status_text FROM `orders`",
			types.DriverPostgres: `SELECT "id", CASE WHEN status = $1 THEN $2 ELSE $3 END AS status_text FROM "orders"`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"'completed'", "'done'", "'pending'"}, bindings)
		assert.NoError(t, err)
	})
}
