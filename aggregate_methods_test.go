package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/iMohamedSheta/xqb/shared/types"

	"github.com/stretchr/testify/assert"
)

func Test_Aggregate(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, binding, err := qb.Select(
			xqb.Sum("price", "total_price"),
			"username",
			"email",
		).Where("id", "=", 15).Limit(1).
			ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT SUM(price) AS total_price, username, email FROM users WHERE id = ? LIMIT 1",
			types.DriverPostgres: "SELECT SUM(price) AS total_price, username, email FROM users WHERE id = $1 LIMIT 1",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{15}, binding)
		assert.NoError(t, err)

	})
}

func Test_DialectExpr(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, binding, err := qb.Select(
			xqb.Sum("price", "total_price"),
			xqb.DateFormat("created_at", "%Y-%m-%d", "created_at"),
			"username",
			"email",
		).Where("id", "=", 15).
			Limit(1).
			ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT SUM(price) AS total_price, DATE_FORMAT(created_at, '%Y-%m-%d') AS created_at, username, email FROM users WHERE id = ? LIMIT 1",
			types.DriverPostgres: "SELECT SUM(price) AS total_price, TO_CHAR(created_at, '%Y-%m-%d') AS created_at, username, email FROM users WHERE id = $1 LIMIT 1",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{15}, binding)
		assert.NoError(t, err)
	})
}

func Test_CountExpression(t *testing.T) {
	expr := xqb.Count("id", "total_users")
	assert.Equal(t, "COUNT(id) AS total_users", expr.SQL)
}

func Test_JsonExtract(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		expr := xqb.JsonExtract("data", "user.name", "username")
		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "JSON_EXTRACT(data, '$.user.name') AS username",
			types.DriverPostgres: "data->'user'->>'name' AS username",
		}
		assert.Equal(t, expectedSql[dialect], expr.Dialects[string(dialect)].SQL)
	})
}

func Test_DateFunctions(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		assert.Equal(t, "DATE(created_at) AS date_only", xqb.Date("created_at", "date_only").SQL)
		dialectExpr := xqb.DateDiff("end_date", "start_date", "diff")

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "DATEDIFF(end_date, start_date) AS diff",
			types.DriverPostgres: "(end_date - start_date) AS diff",
		}
		expr := dialectExpr.Dialects[string(dialect)]

		assert.Equal(t, expectedSql[dialect], expr.SQL)
		assert.Empty(t, expr.Bindings)
	})
}

func Test_StringFunctions(t *testing.T) {
	assert.Equal(t, "LOWER(name) AS lower_name", xqb.Lower("name", "lower_name").SQL)
	assert.Equal(t, "REPLACE(title, 'foo', 'bar') AS updated", xqb.Replace("title", "'foo'", "'bar'", "updated").SQL)
}

func Test_ConcatWithBindings(t *testing.T) {
	expr := xqb.Concat([]string{"first_name", "' '", "last_name"}, "full_name")
	assert.Equal(t, "CONCAT(first_name, ' ', last_name) AS full_name", expr.SQL)
}

func Test_MathExpression(t *testing.T) {
	expr := xqb.Math("price * quantity", "total")
	assert.Equal(t, "price * quantity AS total", expr.SQL)
}

func Test_Coalesce(t *testing.T) {
	expr := xqb.Coalesce([]string{"middle_name", "'N/A'"}, "coalesced_name")
	assert.Equal(t, "COALESCE(middle_name, 'N/A') AS coalesced_name", expr.SQL)
}

func Test_QueryBuilder_Locks(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		sql, b, err := xqb.Table("users").LockForUpdate().ToSQL()
		assert.Equal(t, "SELECT * FROM users FOR UPDATE", sql)
		assert.Equal(t, []any(nil), b)
		assert.NoError(t, err)

		qb := xqb.Table("users").SharedLock()

		sql, b, _ = qb.ToSQL()
		assert.Equal(t, "SELECT * FROM users LOCK IN SHARE MODE", sql)
		assert.Equal(t, []any(nil), b)

		qb = xqb.Table("users").SharedLock().NoWaitLocked()
		sql, b, _ = qb.ToSQL()
		assert.Equal(t, "SELECT * FROM users LOCK IN SHARE MODE NOWAIT", sql)
		assert.Equal(t, []any(nil), b)

		qb = xqb.Table("users").LockForUpdate().SkipLocked()
		sql, b, _ = qb.ToSQL()
		assert.Equal(t, "SELECT * FROM users FOR UPDATE SKIP LOCKED", sql)
		assert.Equal(t, []any(nil), b)

		qb = xqb.Table("users")
		qb.SetDialect(types.DriverPostgres)
		qb.Where("id", "=", 15).LockNoKeyUpdate().SkipLocked()
		sql, b, _ = qb.ToSQL()
		assert.Equal(t, "SELECT * FROM users WHERE id = $1 FOR NO KEY UPDATE SKIP LOCKED", sql)
		assert.Equal(t, []any{15}, b)
	})
}

func Test_Upper(t *testing.T) {
	expr := xqb.Upper("name", "upper_name")
	assert.Equal(t, "UPPER(name) AS upper_name", expr.SQL)
}

func Test_Length(t *testing.T) {
	expr := xqb.Length("description", "desc_length")
	assert.Equal(t, "LENGTH(description) AS desc_length", expr.SQL)
}

func Test_Trim(t *testing.T) {
	expr := xqb.Trim("username", "trimmed_username")
	assert.Equal(t, "TRIM(username) AS trimmed_username", expr.SQL)
}

func Test_DateAdd(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		dialectExpr := xqb.DateAdd("created_at", "7", "DAY", "next_week")
		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "DATE_ADD(created_at, INTERVAL 7 DAY) AS next_week",
			types.DriverPostgres: "created_at + INTERVAL '7 day' AS next_week",
		}

		expr := dialectExpr.Dialects[string(dialect)]

		assert.Equal(t, expectedSql[dialect], expr.SQL)
		assert.Empty(t, expr.Bindings)
	})
}

func Test_DateSub(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		dialectExpr := xqb.DateSub("created_at", "1", "MONTH", "last_month")
		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "DATE_SUB(created_at, INTERVAL 1 MONTH) AS last_month",
			types.DriverPostgres: "created_at - INTERVAL '1 month' AS last_month",
		}

		expr := dialectExpr.Dialects[string(dialect)]
		assert.Equal(t, expectedSql[dialect], expr.SQL)
		assert.Empty(t, expr.Bindings)
	})
}

func Test_Substring(t *testing.T) {
	expr := xqb.Substring("bio", 1, 10, "short_bio")
	assert.Equal(t, "SUBSTRING(bio, 1, 10) AS short_bio", expr.SQL)
}

func Test_QueryBuilder_Upper_Length_Trim(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.Select(
			xqb.Upper("name", "upper_name"),
			xqb.Length("bio", "bio_length"),
			xqb.Trim("username", "trimmed_username"),
		).Where("active", "=", true).ToSQL()
		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT UPPER(name) AS upper_name, LENGTH(bio) AS bio_length, TRIM(username) AS trimmed_username FROM users WHERE active = ?",
			types.DriverPostgres: "SELECT UPPER(name) AS upper_name, LENGTH(bio) AS bio_length, TRIM(username) AS trimmed_username FROM users WHERE active = $1",
		}

		assert.NoError(t, err)
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{true}, bindings)
	})
}

func Test_QueryBuilder_DateAdd_DateSub(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("events").SetDialect(dialect)
		sql, bindings, err := qb.Select(
			xqb.DateAdd("event_date", "1", "DAY", "tomorrow"),
			xqb.DateSub("event_date", "7", "DAY", "last_week"),
		).Where("status", "=", "open").ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT DATE_ADD(event_date, INTERVAL 1 DAY) AS tomorrow, DATE_SUB(event_date, INTERVAL 7 DAY) AS last_week FROM events WHERE status = ?",
			types.DriverPostgres: "SELECT event_date + INTERVAL '1 day' AS tomorrow, event_date - INTERVAL '7 day' AS last_week FROM events WHERE status = $1",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"open"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_QueryBuilder_AggregateMethods(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("test_table").SetDialect(dialect)
		sql, bindings, err := qb.Select(
			xqb.Count("id", "total_count"),
			xqb.Sum("amount", "total_amount"),
			xqb.Avg("score", "avg_score"),
			xqb.Min("age", "min_age"),
			xqb.Max("salary", "max_salary"),
			xqb.JsonExtract("data", "user.email", "user_email"),
			xqb.Math("price * quantity", "total_price"),
			xqb.Date("created_at", "created_date"),
			xqb.DateDiff("end_date", "start_date", "days_between"),
			xqb.DateAdd("created_at", "1", "DAY", "next_day"),
			xqb.DateSub("created_at", "1", "MONTH", "prev_month"),
			xqb.DateFormat("created_at", "%Y-%m-%d", "formatted_date"),
			xqb.Coalesce([]string{"middle_name", "'N/A'"}, "coalesced_name"),
			xqb.Concat([]string{"first_name", "' '", "last_name"}, "full_name"),
			xqb.Lower("email", "lower_email"),
			xqb.Upper("username", "upper_username"),
			xqb.Length("bio", "bio_length"),
			xqb.Trim("nickname", "trimmed_nickname"),
			xqb.Replace("title", "'foo'", "'bar'", "replaced_title"),
			xqb.Substring("description", 1, 10, "short_desc"),
		).Where("active", "=", true).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL: "SELECT COUNT(id) AS total_count, SUM(amount) AS total_amount, " +
				"AVG(score) AS avg_score, MIN(age) AS min_age, MAX(salary) AS max_salary, JSON_EXTRACT(data, '$.user.email') AS user_email, " +
				"price * quantity AS total_price, DATE(created_at) AS created_date, " +
				"DATEDIFF(end_date, start_date) AS days_between, DATE_ADD(created_at, INTERVAL 1 DAY) AS next_day, DATE_SUB(created_at, INTERVAL 1 MONTH) AS prev_month, " +
				"DATE_FORMAT(created_at, '%Y-%m-%d') AS formatted_date, COALESCE(middle_name, 'N/A') AS coalesced_name, CONCAT(first_name, ' ', last_name) AS full_name, " +
				"LOWER(email) AS lower_email, UPPER(username) AS upper_username, LENGTH(bio) AS bio_length, TRIM(nickname) AS trimmed_nickname, " +
				"REPLACE(title, 'foo', 'bar') AS replaced_title, SUBSTRING(description, 1, 10) AS short_desc FROM test_table WHERE active = ?",
			types.DriverPostgres: "SELECT COUNT(id) AS total_count, SUM(amount) AS total_amount, " +
				"AVG(score) AS avg_score, MIN(age) AS min_age, MAX(salary) AS max_salary, data->'user'->>'email' AS user_email, " +
				"price * quantity AS total_price, DATE(created_at) AS created_date, " +
				"(end_date - start_date) AS days_between, created_at + INTERVAL '1 day' AS next_day, created_at - INTERVAL '1 month' AS prev_month, " +
				"TO_CHAR(created_at, '%Y-%m-%d') AS formatted_date, COALESCE(middle_name, 'N/A') AS coalesced_name, CONCAT(first_name, ' ', last_name) AS full_name, " +
				"LOWER(email) AS lower_email, UPPER(username) AS upper_username, LENGTH(bio) AS bio_length, TRIM(nickname) AS trimmed_nickname, " +
				"REPLACE(title, 'foo', 'bar') AS replaced_title, SUBSTRING(description, 1, 10) AS short_desc FROM test_table WHERE active = $1",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{true}, bindings)
		assert.NoError(t, err)

		dialectExprDateFormat := xqb.DateFormat("created_at", "%Y-%m-%d", "formatted_date")
		expectedSql = map[types.Driver]string{
			types.DriverMySQL:    "DATE_FORMAT(created_at, '%Y-%m-%d') AS formatted_date",
			types.DriverPostgres: "TO_CHAR(created_at, '%Y-%m-%d') AS formatted_date",
		}

		assert.Equal(t, expectedSql[dialect], dialectExprDateFormat.Dialects[string(dialect)].SQL)
		assert.Empty(t, dialectExprDateFormat.Dialects[string(dialect)].Bindings)
	})
}

func Test_QueryBuilder_AggregateMethods_2(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("coverage_table").SetDialect(dialect)
		sql, bindings, err := qb.Select(
			xqb.Count("*", ""),
			xqb.Count("id", "cnt"),
			xqb.Sum("amount", ""),
			xqb.Sum("amount", "total_amount"),
			xqb.Avg("score", ""),
			xqb.Avg("score", "avg_score"),
			xqb.Min("age", ""),
			xqb.Min("age", "min_age"),
			xqb.Max("salary", ""),
			xqb.Max("salary", "max_salary"),
			xqb.JsonExtract("data", "user.name", ""),
			xqb.JsonExtract("data", "user.name", "user_name"),
			xqb.Math("price * quantity", ""),
			xqb.Math("price * quantity + tax", "total_price"),
			xqb.Date("created_at", ""),
			xqb.Date("created_at", "created_date"),
			xqb.DateDiff("end_date", "start_date", ""),
			xqb.DateDiff("end_date", "start_date", "days_between"),
			xqb.DateAdd("created_at", "1", "DAY", ""),
			xqb.DateAdd("created_at", "1", "DAY", "next_day"),
			xqb.DateSub("created_at", "1", "MONTH", ""),
			xqb.DateSub("created_at", "1", "MONTH", "prev_month"),
			xqb.DateFormat("created_at", "%Y-%m-%d", ""),
			xqb.DateFormat("created_at", "%Y-%m-%d", "formatted_date"),
			xqb.Coalesce([]string{"middle_name", "'N/A'"}, ""),
			xqb.Coalesce([]string{"middle_name", "'N/A'"}, "coalesced_name"),
			xqb.Concat([]string{"first_name", "' '", "last_name"}, ""),
			xqb.Concat([]string{"first_name", "' '", "last_name"}, "full_name"),
			xqb.Lower("email", ""),
			xqb.Lower("email", "lower_email"),
			xqb.Upper("username", ""),
			xqb.Upper("username", "upper_username"),
			xqb.Length("bio", ""),
			xqb.Length("bio", "bio_length"),
			xqb.Trim("nickname", ""),
			xqb.Trim("nickname", "trimmed_nickname"),
			xqb.Replace("title", "'foo'", "'bar'", ""),
			xqb.Replace("title", "'foo'", "'bar'", "replaced_title"),
			xqb.Substring("description", 1, 10, ""),
			xqb.Substring("description", 1, 10, "short_desc"),
		).Where(
			xqb.Lower("status", ""), "=", "active",
		).GroupBy(
			xqb.Date("created_at", ""),
			xqb.Upper("region", ""),
		).Having(
			xqb.Sum("amount", ""), ">", 1000,
		).OrderBy(
			xqb.Length("bio", ""), "DESC",
		).Limit(5).Offset(10).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL: "SELECT COUNT(*), COUNT(id) AS cnt, SUM(amount), SUM(amount) AS total_amount, AVG(score), AVG(score) AS avg_score, MIN(age), " +
				"MIN(age) AS min_age, MAX(salary), MAX(salary) AS max_salary, JSON_EXTRACT(data, '$.user.name'), JSON_EXTRACT(data, '$.user.name') AS user_name, price * quantity, " +
				"price * quantity + tax AS total_price, DATE(created_at), DATE(created_at) AS created_date, DATEDIFF(end_date, start_date), DATEDIFF(end_date, start_date) AS days_between, " +
				"DATE_ADD(created_at, INTERVAL 1 DAY), DATE_ADD(created_at, INTERVAL 1 DAY) AS next_day, DATE_SUB(created_at, INTERVAL 1 MONTH), DATE_SUB(created_at, INTERVAL 1 MONTH) AS prev_month, " +
				"DATE_FORMAT(created_at, '%Y-%m-%d'), DATE_FORMAT(created_at, '%Y-%m-%d') AS formatted_date, COALESCE(middle_name, 'N/A'), COALESCE(middle_name, 'N/A') AS coalesced_name, CONCAT(first_name, ' ', last_name), " +
				"CONCAT(first_name, ' ', last_name) AS full_name, LOWER(email), LOWER(email) AS lower_email, UPPER(username), UPPER(username) AS upper_username, LENGTH(bio), LENGTH(bio) AS bio_length, " +
				"TRIM(nickname), TRIM(nickname) AS trimmed_nickname, REPLACE(title, 'foo', 'bar'), REPLACE(title, 'foo', 'bar') AS replaced_title, SUBSTRING(description, 1, 10), " +
				"SUBSTRING(description, 1, 10) AS short_desc FROM coverage_table WHERE LOWER(status) = ? GROUP BY DATE(created_at), UPPER(region) HAVING SUM(amount) > ? ORDER BY LENGTH(bio) DESC LIMIT 5 OFFSET 10",
			types.DriverPostgres: "SELECT COUNT(*), COUNT(id) AS cnt, SUM(amount), SUM(amount) AS total_amount, AVG(score), AVG(score) AS avg_score, MIN(age), " +
				"MIN(age) AS min_age, MAX(salary), MAX(salary) AS max_salary, data->'user'->>'name', data->'user'->>'name' AS user_name, price * quantity, " +
				"price * quantity + tax AS total_price, DATE(created_at), DATE(created_at) AS created_date, (end_date - start_date), (end_date - start_date) AS days_between, " +
				"created_at + INTERVAL '1 day', created_at + INTERVAL '1 day' AS next_day, created_at - INTERVAL '1 month', created_at - INTERVAL '1 month' AS prev_month, " +
				"TO_CHAR(created_at, '%Y-%m-%d'), TO_CHAR(created_at, '%Y-%m-%d') AS formatted_date, COALESCE(middle_name, 'N/A'), COALESCE(middle_name, 'N/A') AS coalesced_name, CONCAT(first_name, ' ', last_name), " +
				"CONCAT(first_name, ' ', last_name) AS full_name, LOWER(email), LOWER(email) AS lower_email, UPPER(username), UPPER(username) AS upper_username, LENGTH(bio), LENGTH(bio) AS bio_length, " +
				"TRIM(nickname), TRIM(nickname) AS trimmed_nickname, REPLACE(title, 'foo', 'bar'), REPLACE(title, 'foo', 'bar') AS replaced_title, SUBSTRING(description, 1, 10), " +
				"SUBSTRING(description, 1, 10) AS short_desc FROM coverage_table WHERE LOWER(status) = $1 GROUP BY DATE(created_at), UPPER(region) HAVING SUM(amount) > $2 ORDER BY LENGTH(bio) DESC LIMIT 5 OFFSET 10",
		}

		assert.NoError(t, err)
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"active", 1000}, bindings)
	})
}
